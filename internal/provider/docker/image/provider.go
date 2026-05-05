package image

import (
	"context"

	mobyclient "github.com/moby/moby/client"
)

type dockerClient interface {
	ImagePull(ctx context.Context, refStr string, options mobyclient.ImagePullOptions) (mobyclient.ImagePullResponse, error)
	ImagePush(ctx context.Context, image string, options mobyclient.ImagePushOptions) (mobyclient.ImagePushResponse, error)
	ImageRemove(ctx context.Context, imageID string, options mobyclient.ImageRemoveOptions) (mobyclient.ImageRemoveResult, error)
	ImageTag(ctx context.Context, options mobyclient.ImageTagOptions) (mobyclient.ImageTagResult, error)
	ImageList(ctx context.Context, options mobyclient.ImageListOptions) (mobyclient.ImageListResult, error)
	ImageInspect(ctx context.Context, imageID string, inspectOpts ...mobyclient.ImageInspectOption) (mobyclient.ImageInspectResult, error)
}

type Provider struct {
	client      dockerClient
	withTimeout func(context.Context) (context.Context, context.CancelFunc)
}

func NewProvider(
	client dockerClient,
	withTimeout func(context.Context) (context.Context, context.CancelFunc),
) *Provider {
	return &Provider{
		client:      client,
		withTimeout: withTimeout,
	}
}

func (p *Provider) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return p.withTimeout(ctx)
}
