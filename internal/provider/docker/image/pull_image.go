package image

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func (p *Provider) PullImage(ctx context.Context, params providers.PullImageParams) error {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()

	opts := client.ImagePullOptions{
		All: params.All,
	}

	if params.Platform != nil {
		opts.Platforms = []ocispec.Platform{
			*params.Platform,
		}
	}

	response, err := p.client.ImagePull(ctx, params.Ref, opts)
	if err != nil {
		return fmt.Errorf("pull docker image: %w", err)
	}

	err = response.Wait(ctx)
	closeErr := response.Close()

	if err != nil {
		return fmt.Errorf("wait for docker image pull: %w", err)
	}

	if closeErr != nil {
		return fmt.Errorf("close pull response: %w", closeErr)
	}

	return nil
}
