package container

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) StartContainer(ctx context.Context, params providers.StartContainerParams) error {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()
	_, err := p.client.ContainerStart(ctx, params.Name, client.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("start docker container: %w", err)
	}

	return nil
}
