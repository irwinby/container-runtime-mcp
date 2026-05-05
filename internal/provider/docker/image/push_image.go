package image

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) PushImage(ctx context.Context, params providers.PushImageParams) error {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()

	opts := client.ImagePushOptions{
		All: params.All,
	}

	if params.Platform != nil {
		opts.Platform = params.Platform
	}

	response, err := p.client.ImagePush(ctx, params.Ref, opts)
	if err != nil {
		return fmt.Errorf("push docker image: %w", err)
	}

	err = response.Wait(ctx)
	closeErr := response.Close()

	if err != nil {
		return fmt.Errorf("wait for docker image push: %w", err)
	}

	if closeErr != nil {
		return fmt.Errorf("close push response: %w", closeErr)
	}

	return nil
}
