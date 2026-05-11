package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/irwinby/container-runtime-mcp/internal/api/mcp"
	"github.com/irwinby/container-runtime-mcp/internal/api/mcp/handler/container"
	"github.com/irwinby/container-runtime-mcp/internal/api/mcp/handler/image"
	"github.com/irwinby/container-runtime-mcp/internal/api/mcp/handler/system"
	"github.com/irwinby/container-runtime-mcp/internal/api/mcp/handler/volume"
	"github.com/irwinby/container-runtime-mcp/internal/api/telemetry"
	"github.com/irwinby/container-runtime-mcp/internal/api/telemetry/handler/probe"
	"github.com/irwinby/container-runtime-mcp/internal/config"
	provider "github.com/irwinby/container-runtime-mcp/internal/provider/docker"
	services "github.com/irwinby/container-runtime-mcp/internal/service"
	containerservice "github.com/irwinby/container-runtime-mcp/internal/service/container"
	imageservice "github.com/irwinby/container-runtime-mcp/internal/service/image"
	systemservice "github.com/irwinby/container-runtime-mcp/internal/service/system"
	volumeservice "github.com/irwinby/container-runtime-mcp/internal/service/volume"
	"github.com/irwinby/container-runtime-mcp/pkg/logger"
	"go.uber.org/zap"
)

const shutdownTimeout = 5 * time.Second

func Run(ctx context.Context) error {
	cfg, err := config.LoadFromEnv(ctx)
	if err != nil {
		return fmt.Errorf("load config from environment variables: %w", err)
	}

	return RunWithConfig(ctx, cfg)
}

func RunWithConfig(ctx context.Context, cfg *config.Config) (runErr error) {
	log, err := logger.New(logger.WithLevel(cfg.LogLevel))
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}

	defer func() {
		_ = log.Sync()
	}()

	provider, err := provider.NewProvider(ctx, cfg.RemoteOperationTimeout)
	if err != nil {
		return fmt.Errorf("create provider: %w", err)
	}

	defer func() {
		err := provider.Close()
		if err != nil && runErr == nil {
			runErr = fmt.Errorf("close provider: %w", err)
		}
	}()

	policy := services.NewPolicy(cfg.ReadOnly)

	containerService := containerservice.NewService(provider, policy, log)
	imageService := imageservice.NewService(provider, policy, log)
	volumeService := volumeservice.NewService(provider, policy, log)
	systemService := systemservice.NewService(provider, log)

	log.Info("starting mcp server",
		zap.String("transport", string(cfg.TransportConfig.Type)),
		zap.String("name", cfg.Name),
		zap.String("version", cfg.Version),
	)

	mcpServer, mcpErrs, err := startMCPServer(ctx, cfg, containerService, imageService, volumeService, systemService)
	if err != nil {
		return err
	}

	telemetryServer, telemetryErrs, err := maybeStartTelemetry(ctx, cfg, log, systemService, mcpErrs)
	if err != nil {
		return err
	}

	return waitAndShutdown(ctx, log, mcpServer, telemetryServer, mcpErrs, telemetryErrs)
}

func startMCPServer(
	ctx context.Context,
	cfg *config.Config,
	containerService *containerservice.Service,
	imageService *imageservice.Service,
	volumeService *volumeservice.Service,
	systemService *systemservice.Service,
) (mcp.Server, chan error, error) {
	server, err := mcp.NewServer(
		cfg.MCPServer,
		mcp.NewHandlers(
			container.NewToolsHandler(containerService),
			image.NewToolsHandler(imageService),
			volume.NewToolsHandler(volumeService),
			system.NewToolsHandler(systemService),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create mcp server: %w", err)
	}

	errs := make(chan error, 1)

	go func() {
		errs <- server.Run(ctx)
	}()

	return server, errs, nil
}

func maybeStartTelemetry(
	ctx context.Context,
	cfg *config.Config,
	log *zap.Logger,
	systemService *systemservice.Service,
	mcpErrs chan error,
) (*telemetry.Server, chan error, error) {
	if !cfg.Telemetry.Enabled {
		return nil, nil, nil
	}

	telemetryServer, err := telemetry.NewServer(
		cfg.Telemetry,
		telemetry.NewHandler(
			probe.NewHandler(systemService),
		),
	)
	if err != nil {
		go func() { <-mcpErrs }()
		return nil, nil, fmt.Errorf("create telemetry server: %w", err)
	}

	log.Info("starting telemetry server",
		zap.String("addr", cfg.Telemetry.Addr),
		zap.Bool("pprof", cfg.Telemetry.PPROFEnabled),
	)

	errs := make(chan error, 1)

	go func() {
		errs <- telemetryServer.Run(ctx)
	}()

	return telemetryServer, errs, nil
}

func waitAndShutdown(
	ctx context.Context,
	log *zap.Logger,
	server mcp.Server,
	telemetryServer *telemetry.Server,
	mcpErrs chan error,
	telemetryErrs chan error,
) error {
	var firstErr error

	select {
	case <-ctx.Done():
		log.Info("shutting down servers")
	case err := <-mcpErrs:
		if err != nil {
			firstErr = fmt.Errorf("run mcp server: %w", err)
		}

		log.Info("mcp server exited, shutting down")
	case err := <-telemetryErrs:
		if err != nil {
			firstErr = fmt.Errorf("run telemetry server: %w", err)
		}

		log.Info("telemetry server exited, shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		firstErr = setFirstError(firstErr, fmt.Errorf("shutdown mcp server: %w", err))
	}

	if telemetryServer != nil {
		if err := telemetryServer.Shutdown(shutdownCtx); err != nil {
			firstErr = setFirstError(firstErr, fmt.Errorf("shutdown telemetry server: %w", err))
		}
	}

	if err := <-mcpErrs; err != nil && !errors.Is(err, http.ErrServerClosed) {
		firstErr = setFirstError(firstErr, fmt.Errorf("run mcp server: %w", err))
	}

	if telemetryErrs != nil {
		if err := <-telemetryErrs; err != nil && !errors.Is(err, http.ErrServerClosed) {
			firstErr = setFirstError(firstErr, fmt.Errorf("run telemetry server: %w", err))
		}
	}

	log.Info("servers stopped")

	return firstErr
}

func setFirstError(current, err error) error {
	if current == nil {
		return err
	}

	return current
}
