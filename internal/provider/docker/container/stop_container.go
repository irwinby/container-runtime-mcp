package container

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) StopContainer(ctx context.Context, params providers.StopContainerParams) error {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()

	opts := client.ContainerStopOptions{}

	if params.Signal != "" {
		opts.Signal = params.Signal
	}

	if params.TimeoutSeconds != nil {
		opts.Timeout = params.TimeoutSeconds
	}

	_, err := p.client.ContainerStop(ctx, params.Name, opts)
	if err != nil {
		return fmt.Errorf("stop docker container: %w", err)
	}

	return nil
}
