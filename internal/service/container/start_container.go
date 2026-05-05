package container

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// StartContainerParams holds the parameters for starting a container.
type StartContainerParams struct {
	Name string
}

func NewStartContainerParams() StartContainerParams {
	return StartContainerParams{}
}

func (p StartContainerParams) SetName(name string) StartContainerParams {
	p.Name = strings.TrimSpace(name)
	return p
}

// Validate checks that required fields are present and trims whitespace.
func (p StartContainerParams) Validate() (StartContainerParams, error) {
	if p.Name == "" {
		return StartContainerParams{}, fmt.Errorf("name is required")
	}

	return p, nil
}

func (s *Service) StartContainer(ctx context.Context, params StartContainerParams) error {
	err := s.policy.IsWriteAllowed()
	if err != nil {
		s.logger.Warn("start container blocked by policy", zap.Error(err))
		return fmt.Errorf("check if write is allowed: %w", err)
	}

	params, err = params.Validate()
	if err != nil {
		s.logger.Warn("start container validation failed", zap.Error(err))
		return fmt.Errorf("validate start container params: %w", err)
	}

	s.logger.Info("starting container", zap.String("name", params.Name))

	err = s.providerClient.StartContainer(ctx, providers.StartContainerParams{
		Name: params.Name,
	})
	if err != nil {
		s.logger.Error("start container failed", zap.String("name", params.Name), zap.Error(err))
		return fmt.Errorf("start container: %w", err)
	}

	s.logger.Info("container started", zap.String("name", params.Name))

	return nil
}
