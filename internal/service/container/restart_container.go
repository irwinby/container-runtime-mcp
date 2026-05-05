package container

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// RestartContainerParams holds the parameters for restarting a container.
type RestartContainerParams struct {
	Name           string
	Signal         string
	TimeoutSeconds *int
}

func NewRestartContainerParams() RestartContainerParams {
	return RestartContainerParams{}
}

func (p RestartContainerParams) SetName(name string) RestartContainerParams {
	p.Name = strings.TrimSpace(name)
	return p
}

func (p RestartContainerParams) SetSignal(signal string) RestartContainerParams {
	p.Signal = strings.TrimSpace(signal)
	return p
}

func (p RestartContainerParams) SetTimeoutSeconds(timeoutSeconds *int) RestartContainerParams {
	p.TimeoutSeconds = timeoutSeconds
	return p
}

// Validate checks that required fields are present and trims whitespace.
func (p RestartContainerParams) Validate() (RestartContainerParams, error) {
	if p.Name == "" {
		return RestartContainerParams{}, fmt.Errorf("name is required")
	}

	return p, nil
}

func (s *Service) RestartContainer(ctx context.Context, params RestartContainerParams) error {
	err := s.policy.IsWriteAllowed()
	if err != nil {
		s.logger.Warn("restart container blocked by policy", zap.Error(err))
		return fmt.Errorf("check if write is allowed: %w", err)
	}

	params, err = params.Validate()
	if err != nil {
		s.logger.Warn("restart container validation failed", zap.Error(err))
		return fmt.Errorf("validate restart container params: %w", err)
	}

	s.logger.Info("restarting container", zap.String("name", params.Name))

	err = s.providerClient.RestartContainer(ctx, providers.RestartContainerParams{
		Name:           params.Name,
		Signal:         params.Signal,
		TimeoutSeconds: params.TimeoutSeconds,
	})
	if err != nil {
		s.logger.Error("restart container failed", zap.String("name", params.Name), zap.Error(err))
		return fmt.Errorf("restart container: %w", err)
	}

	s.logger.Info("container restarted", zap.String("name", params.Name))

	return nil
}
