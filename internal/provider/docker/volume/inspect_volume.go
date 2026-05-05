package volume

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) InspectVolume(ctx context.Context, params providers.InspectVolumeParams) (providers.VolumeInspect, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()

	result, err := p.client.VolumeInspect(ctx, params.Name, client.VolumeInspectOptions{})
	if err != nil {
		return providers.VolumeInspect{}, fmt.Errorf("inspect docker volume: %w", err)
	}

	return providers.NewVolumeInspect().
		SetName(result.Volume.Name).
		SetDriver(result.Volume.Driver).
		SetMountpoint(result.Volume.Mountpoint).
		SetLabels(result.Volume.Labels).
		SetScope(result.Volume.Scope), nil
}
