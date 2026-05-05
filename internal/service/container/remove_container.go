package container

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// RemoveContainerParams holds the parameters for removing a container.
type RemoveContainerParams struct {
	Name          string
	Force         bool
	RemoveVolumes bool
	RemoveLinks   bool
}

func NewRemoveContainerParams() RemoveContainerParams {
	return RemoveContainerParams{}
}

func (p RemoveContainerParams) SetName(name string) RemoveContainerParams {
	p.Name = strings.TrimSpace(name)

	return p
}

func (p RemoveContainerParams) SetForce(force bool) RemoveContainerParams {
	p.Force = force

	return p
}

func (p RemoveContainerParams) SetRemoveVolumes(removeVolumes bool) RemoveContainerParams {
	p.RemoveVolumes = removeVolumes

	return p
}

func (p RemoveContainerParams) SetRemoveLinks(removeLinks bool) RemoveContainerParams {
	p.RemoveLinks = removeLinks

	return p
}

// Validate checks that required fields are present and trims whitespace.
func (i RemoveContainerParams) Validate() (RemoveContainerParams, error) {
	if i.Name == "" {
		return RemoveContainerParams{}, fmt.Errorf("name is required")
	}

	return i, nil
}

func (s *Service) RemoveContainer(ctx context.Context, params RemoveContainerParams) error {
	err := s.policy.IsWriteAllowed()
	if err != nil {
		s.logger.Warn("remove container blocked by policy", zap.Error(err))
		return fmt.Errorf("check if write is allowed: %w", err)
	}

	params, err = params.Validate()
	if err != nil {
		s.logger.Warn("remove container validation failed", zap.Error(err))
		return fmt.Errorf("validate remove container params: %w", err)
	}

	s.logger.Info("removing container", zap.String("name", params.Name))

	err = s.providerClient.RemoveContainer(ctx, providers.RemoveContainerParams{
		Name:          params.Name,
		Force:         params.Force,
		RemoveVolumes: params.RemoveVolumes,
		RemoveLinks:   params.RemoveLinks,
	})
	if err != nil {
		s.logger.Error("remove container failed", zap.String("name", params.Name), zap.Error(err))
		return fmt.Errorf("remove container: %w", err)
	}

	s.logger.Info("container removed", zap.String("name", params.Name))

	return nil
}
