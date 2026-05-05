package system

import (
	"context"

	mobyclient "github.com/moby/moby/client"
)

type dockerClient interface {
	Ping(ctx context.Context, options mobyclient.PingOptions) (mobyclient.PingResult, error)
	Info(ctx context.Context, options mobyclient.InfoOptions) (mobyclient.SystemInfoResult, error)
	ServerVersion(ctx context.Context, options mobyclient.ServerVersionOptions) (mobyclient.ServerVersionResult, error)
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
