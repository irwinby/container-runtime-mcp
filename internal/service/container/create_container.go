package container

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// CreateContainerParams holds the parameters for creating a container.
type CreateContainerParams struct {
	Name  string
	Image string
}

func NewCreateContainerParams() CreateContainerParams {
	return CreateContainerParams{}
}

func (p CreateContainerParams) SetName(name string) CreateContainerParams {
	p.Name = strings.TrimSpace(name)

	return p
}

func (p CreateContainerParams) SetImage(image string) CreateContainerParams {
	p.Image = strings.TrimSpace(image)

	return p
}

// Validate checks that required fields are present and trims whitespace.
func (p CreateContainerParams) Validate() (CreateContainerParams, error) {
	if p.Name == "" {
		return CreateContainerParams{}, fmt.Errorf("name is required")
	}

	if p.Image == "" {
		return CreateContainerParams{}, fmt.Errorf("image is required")
	}

	return p, nil
}

func (s *Service) CreateContainer(ctx context.Context, params CreateContainerParams) (string, error) {
	err := s.policy.IsWriteAllowed()
	if err != nil {
		s.logger.Warn("create container blocked by policy", zap.Error(err))
		return "", fmt.Errorf("check if write is allowed: %w", err)
	}

	params, err = params.Validate()
	if err != nil {
		s.logger.Warn("create container validation failed", zap.Error(err))
		return "", fmt.Errorf("validate create container params: %w", err)
	}

	s.logger.Info("creating container", zap.String("name", params.Name), zap.String("image", params.Image))

	id, err := s.providerClient.CreateContainer(ctx, providers.CreateContainerParams{
		Name:  params.Name,
		Image: params.Image,
	})
	if err != nil {
		s.logger.Error("create container failed", zap.String("name", params.Name), zap.Error(err))
		return "", fmt.Errorf("create container: %w", err)
	}

	s.logger.Info("container created", zap.String("id", id))

	return id, nil
}
