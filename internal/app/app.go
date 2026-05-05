package app

import (
	"context"
	"fmt"
	"time"

	"github.com/irwinby/container-runtime-mcp/internal/api/mcp"
	"github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/container"
	"github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/image"
	"github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/system"
	"github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/volume"
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
		return fmt.Errorf("create mcp server: %w", err)
	}

	errs := make(chan error, 1)

	go func() {
		errs <- server.Run(ctx)
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down mcp server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			go func() { <-errs }()

			return fmt.Errorf("shutdown mcp server: %w", err)
		}

		err = <-errs
		if err != nil {
			return fmt.Errorf("run mcp server: %w", err)
		}

		log.Info("mcp server stopped")

		return nil
	case err := <-errs:
		if err != nil {
			return fmt.Errorf("run mcp server: %w", err)
		}

		return nil
	}
}
