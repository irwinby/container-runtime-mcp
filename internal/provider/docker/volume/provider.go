package volume

import (
	"context"

	mobyclient "github.com/moby/moby/client"
)

type dockerClient interface {
	VolumeList(ctx context.Context, options mobyclient.VolumeListOptions) (mobyclient.VolumeListResult, error)
	VolumeInspect(ctx context.Context, volumeID string, options mobyclient.VolumeInspectOptions) (mobyclient.VolumeInspectResult, error)
	VolumeCreate(ctx context.Context, options mobyclient.VolumeCreateOptions) (mobyclient.VolumeCreateResult, error)
	VolumeRemove(ctx context.Context, volumeID string, options mobyclient.VolumeRemoveOptions) (mobyclient.VolumeRemoveResult, error)
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
