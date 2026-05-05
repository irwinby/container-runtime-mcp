package image

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func (p *Provider) RemoveImage(ctx context.Context, params providers.RemoveImageParams) error {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()
	opts := client.ImageRemoveOptions{
		Force:         params.Force,
		PruneChildren: params.PruneChildren,
	}

	if params.Platform != nil {
		opts.Platforms = []ocispec.Platform{*params.Platform}
	}

	_, err := p.client.ImageRemove(ctx, params.Ref, opts)
	if err != nil {
		return fmt.Errorf("remove docker image: %w", err)
	}

	return nil
}
