package volume

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// ListVolumesParams holds the parameters for listing volumes.
type ListVolumesParams struct {
	Dangling bool
}

func NewListVolumesParams() ListVolumesParams {
	return ListVolumesParams{}
}

func (p ListVolumesParams) SetDangling(dangling bool) ListVolumesParams {
	p.Dangling = dangling
	return p
}

func (s *Service) ListVolumes(ctx context.Context, params ListVolumesParams) ([]Volume, error) {
	s.logger.Debug("listing volumes")

	result, err := s.providerClient.ListVolumes(ctx, providers.ListVolumesParams{
		Dangling: params.Dangling,
	})
	if err != nil {
		s.logger.Error("list volumes failed", zap.Error(err))
		return nil, fmt.Errorf("list volumes: %w", err)
	}

	volumes := make([]Volume, 0, len(result))

	for _, volume := range result {
		volumes = append(volumes, Volume{
			Name:       volume.Name,
			Driver:     volume.Driver,
			Mountpoint: volume.Mountpoint,
			Labels:     volume.Labels,
			Scope:      volume.Scope,
		})
	}

	return volumes, nil
}
