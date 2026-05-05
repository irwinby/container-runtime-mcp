package container

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) InspectContainer(ctx context.Context, params providers.InspectContainerParams) (providers.ContainerInspect, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()
	result, err := p.client.ContainerInspect(ctx, params.Name, client.ContainerInspectOptions{})
	if err != nil {
		return providers.ContainerInspect{}, fmt.Errorf("inspect docker container: %w", err)
	}

	info := providers.NewContainerInspect().
		SetID(result.Container.ID).
		SetName(result.Container.Name).
		SetImage(result.Container.Image).
		SetPath(result.Container.Path).
		SetArgs(result.Container.Args).
		SetRestartCount(result.Container.RestartCount).
		SetCreated(result.Container.Created)

	if result.Container.State != nil {
		info = info.
			SetState(string(result.Container.State.Status)).
			SetStatus(string(result.Container.State.Status))
	}

	return info, nil
}
