package container

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// StopContainerParams holds the parameters for stopping a container.
type StopContainerParams struct {
	Name           string
	Signal         string
	TimeoutSeconds *int
}

func NewStopContainerParams() StopContainerParams {
	return StopContainerParams{}
}

func (p StopContainerParams) SetName(name string) StopContainerParams {
	p.Name = strings.TrimSpace(name)
	return p
}

func (p StopContainerParams) SetSignal(signal string) StopContainerParams {
	p.Signal = strings.TrimSpace(signal)
	return p
}

func (p StopContainerParams) SetTimeoutSeconds(timeoutSeconds *int) StopContainerParams {
	p.TimeoutSeconds = timeoutSeconds
	return p
}

// Validate checks that required fields are present and trims whitespace.
func (p StopContainerParams) Validate() (StopContainerParams, error) {
	if p.Name == "" {
		return StopContainerParams{}, fmt.Errorf("name is required")
	}

	return p, nil
}

func (s *Service) StopContainer(ctx context.Context, params StopContainerParams) error {
	err := s.policy.IsWriteAllowed()
	if err != nil {
		s.logger.Warn("stop container blocked by policy", zap.Error(err))
		return fmt.Errorf("check if write is allowed: %w", err)
	}

	params, err = params.Validate()
	if err != nil {
		s.logger.Warn("stop container validation failed", zap.Error(err))
		return fmt.Errorf("validate stop container params: %w", err)
	}

	s.logger.Info("stopping container", zap.String("name", params.Name))

	err = s.providerClient.StopContainer(ctx, providers.StopContainerParams{
		Name:           params.Name,
		Signal:         params.Signal,
		TimeoutSeconds: params.TimeoutSeconds,
	})
	if err != nil {
		s.logger.Error("stop container failed", zap.String("name", params.Name), zap.Error(err))
		return fmt.Errorf("stop container: %w", err)
	}

	s.logger.Info("container stopped", zap.String("name", params.Name))

	return nil
}
