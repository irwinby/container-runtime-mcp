package container

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) ListContainers(ctx context.Context, params providers.ListContainersParams) ([]providers.Container, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()

	result, err := p.client.ContainerList(ctx, client.ContainerListOptions{
		All:    params.All,
		Limit:  params.Limit,
		Size:   params.Size,
		Latest: params.Latest,
	})
	if err != nil {
		return nil, fmt.Errorf("list docker containers: %w", err)
	}

	containers := make([]providers.Container, 0, len(result.Items))

	for _, item := range result.Items {
		containers = append(containers, providers.NewContainer().
			SetID(item.ID).
			SetNames(item.Names).
			SetImage(item.Image).
			SetState(string(item.State)).
			SetStatus(item.Status).
			SetCreated(item.Created))
	}

	return containers, nil
}
