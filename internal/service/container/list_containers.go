package container

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// ListContainersParams holds the parameters for listing containers.
type ListContainersParams struct {
	All    bool
	Limit  int
	Size   bool
	Latest bool
}

func NewListContainersParams() ListContainersParams {
	return ListContainersParams{}
}

func (p ListContainersParams) SetAll(all bool) ListContainersParams {
	p.All = all
	return p
}

func (p ListContainersParams) SetLimit(limit int) ListContainersParams {
	p.Limit = limit
	return p
}

func (p ListContainersParams) SetSize(size bool) ListContainersParams {
	p.Size = size
	return p
}

func (p ListContainersParams) SetLatest(latest bool) ListContainersParams {
	p.Latest = latest
	return p
}

// Validate checks that required fields are present.
func (p ListContainersParams) Validate() (ListContainersParams, error) {
	if p.Limit < 0 {
		return ListContainersParams{}, fmt.Errorf("limit must be non-negative")
	}

	return p, nil
}

func (s *Service) ListContainers(ctx context.Context, params ListContainersParams) ([]Container, error) {
	params, err := params.Validate()
	if err != nil {
		s.logger.Warn("list containers validation failed", zap.Error(err))
		return nil, fmt.Errorf("validate list containers params: %w", err)
	}

	s.logger.Debug("listing containers")

	containers, err := s.providerClient.ListContainers(ctx, providers.ListContainersParams{
		All:    params.All,
		Limit:  params.Limit,
		Size:   params.Size,
		Latest: params.Latest,
	})
	if err != nil {
		s.logger.Error("list containers failed", zap.Error(err))
		return nil, fmt.Errorf("list containers: %w", err)
	}

	result := make([]Container, 0, len(containers))

	for _, container := range containers {
		result = append(result, Container{
			ID:      container.ID,
			Names:   container.Names,
			Image:   container.Image,
			State:   container.State,
			Status:  container.Status,
			Created: container.Created,
		})
	}

	return result, nil
}
