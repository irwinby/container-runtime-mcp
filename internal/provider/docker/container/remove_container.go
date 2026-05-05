package container

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) RemoveContainer(ctx context.Context, params providers.RemoveContainerParams) error {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()
	_, err := p.client.ContainerRemove(ctx, params.Name, client.ContainerRemoveOptions{
		Force:         params.Force,
		RemoveVolumes: params.RemoveVolumes,
		RemoveLinks:   params.RemoveLinks,
	})
	if err != nil {
		return fmt.Errorf("remove docker container: %w", err)
	}

	return nil
}
