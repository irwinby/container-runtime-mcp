package volume

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) ListVolumes(ctx context.Context, params providers.ListVolumesParams) ([]providers.Volume, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()

	filters := client.Filters{}
	if params.Dangling {
		filters = filters.Add("dangling", "true")
	}

	result, err := p.client.VolumeList(ctx, client.VolumeListOptions{
		Filters: filters,
	})
	if err != nil {
		return nil, fmt.Errorf("list docker volumes: %w", err)
	}

	volumes := make([]providers.Volume, 0, len(result.Items))

	for _, v := range result.Items {
		volumes = append(volumes, providers.NewVolume().
			SetName(v.Name).
			SetDriver(v.Driver).
			SetMountpoint(v.Mountpoint).
			SetLabels(v.Labels).
			SetScope(v.Scope))
	}

	return volumes, nil
}
