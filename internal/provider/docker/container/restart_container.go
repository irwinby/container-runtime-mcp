package container

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) RestartContainer(ctx context.Context, params providers.RestartContainerParams) error {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()
	opts := client.ContainerRestartOptions{}

	if params.Signal != "" {
		opts.Signal = params.Signal
	}

	if params.TimeoutSeconds != nil {
		opts.Timeout = params.TimeoutSeconds
	}

	_, err := p.client.ContainerRestart(ctx, params.Name, opts)
	if err != nil {
		return fmt.Errorf("restart docker container: %w", err)
	}

	return nil
}
