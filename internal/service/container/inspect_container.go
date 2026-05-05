package container

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// InspectContainerParams holds the parameters for inspecting a container.
type InspectContainerParams struct {
	Name string
}

func NewInspectContainerParams() InspectContainerParams {
	return InspectContainerParams{}
}

func (p InspectContainerParams) SetName(name string) InspectContainerParams {
	p.Name = strings.TrimSpace(name)
	return p
}

// Validate checks that required fields are present and trims whitespace.
func (p InspectContainerParams) Validate() (InspectContainerParams, error) {
	if p.Name == "" {
		return InspectContainerParams{}, fmt.Errorf("name is required")
	}

	return p, nil
}

func (s *Service) InspectContainer(ctx context.Context, params InspectContainerParams) (ContainerInspect, error) {
	params, err := params.Validate()
	if err != nil {
		s.logger.Warn("inspect container validation failed", zap.Error(err))
		return ContainerInspect{}, fmt.Errorf("validate inspect container params: %w", err)
	}

	s.logger.Debug("inspecting container", zap.String("name", params.Name))

	info, err := s.providerClient.InspectContainer(ctx, providers.InspectContainerParams{
		Name: params.Name,
	})
	if err != nil {
		s.logger.Error("inspect container failed", zap.String("name", params.Name), zap.Error(err))
		return ContainerInspect{}, fmt.Errorf("inspect container: %w", err)
	}

	return ContainerInspect{
		ID:           info.ID,
		Name:         info.Name,
		Image:        info.Image,
		State:        info.State,
		Status:       info.Status,
		Created:      info.Created,
		Path:         info.Path,
		Args:         info.Args,
		RestartCount: info.RestartCount,
	}, nil
}
