package image

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) TagImage(ctx context.Context, params providers.TagImageParams) error {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()

	_, err := p.client.ImageTag(ctx, client.ImageTagOptions{
		Source: params.Source,
		Target: params.Target,
	})
	if err != nil {
		return fmt.Errorf("tag docker image: %w", err)
	}

	return nil
}
