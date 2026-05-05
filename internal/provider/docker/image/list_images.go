package image

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) ListImages(ctx context.Context, params providers.ListImagesParams) ([]providers.Image, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()
	result, err := p.client.ImageList(ctx, client.ImageListOptions{
		All:        params.All,
		SharedSize: params.SharedSize,
	})
	if err != nil {
		return nil, fmt.Errorf("list docker images: %w", err)
	}

	images := make([]providers.Image, 0, len(result.Items))
	for _, item := range result.Items {
		images = append(images, providers.NewImage().
			SetID(item.ID).
			SetRepoTags(item.RepoTags).
			SetSize(item.Size).
			SetCreated(item.Created).
			SetContainers(item.Containers))
	}

	return images, nil
}
