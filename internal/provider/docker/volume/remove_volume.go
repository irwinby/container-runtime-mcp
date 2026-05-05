package volume

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) RemoveVolume(ctx context.Context, params providers.RemoveVolumeParams) error {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()

	_, err := p.client.VolumeRemove(ctx, params.Name, client.VolumeRemoveOptions{
		Force: params.Force,
	})
	if err != nil {
		return fmt.Errorf("remove docker volume: %w", err)
	}

	return nil
}
