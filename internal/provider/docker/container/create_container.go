package container

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func (p *Provider) CreateContainer(ctx context.Context, params providers.CreateContainerParams) (string, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()

	result, err := p.client.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name: params.Name,
		Config: &container.Config{
			Image: params.Image,
		},
	})
	if err != nil {
		return "", fmt.Errorf("create docker container: %w", err)
	}

	return result.ID, nil
}
