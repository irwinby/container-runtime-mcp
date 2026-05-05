package volume

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) CreateVolume(ctx context.Context, params providers.CreateVolumeParams) (providers.VolumeInspect, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()

	result, err := p.client.VolumeCreate(ctx, client.VolumeCreateOptions{
		Name:       params.Name,
		Driver:     params.Driver,
		DriverOpts: params.DriverOpts,
		Labels:     params.Labels,
	})
	if err != nil {
		return providers.VolumeInspect{}, fmt.Errorf("create docker volume: %w", err)
	}

	return providers.NewVolumeInspect().
		SetName(result.Volume.Name).
		SetDriver(result.Volume.Driver).
		SetMountpoint(result.Volume.Mountpoint).
		SetLabels(result.Volume.Labels).
		SetScope(result.Volume.Scope), nil
}
