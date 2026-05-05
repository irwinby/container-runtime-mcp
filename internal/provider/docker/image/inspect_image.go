package image

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
)

func (p *Provider) InspectImage(ctx context.Context, params providers.InspectImageParams) (providers.ImageInspect, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()
	result, err := p.client.ImageInspect(ctx, params.Ref)
	if err != nil {
		return providers.ImageInspect{}, fmt.Errorf("inspect docker image: %w", err)
	}

	return providers.NewImageInspect().
		SetID(result.ID).
		SetRepoTags(result.RepoTags).
		SetSize(result.Size).
		SetCreated(result.Created).
		SetArchitecture(result.Architecture).
		SetOS(result.Os), nil
}
