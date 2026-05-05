package container

import (
	"context"

	mobyclient "github.com/moby/moby/client"
)

type dockerClient interface {
	ContainerCreate(ctx context.Context, options mobyclient.ContainerCreateOptions) (mobyclient.ContainerCreateResult, error)
	ContainerStart(ctx context.Context, containerID string, options mobyclient.ContainerStartOptions) (mobyclient.ContainerStartResult, error)
	ContainerStop(ctx context.Context, containerID string, options mobyclient.ContainerStopOptions) (mobyclient.ContainerStopResult, error)
	ContainerRestart(ctx context.Context, containerID string, options mobyclient.ContainerRestartOptions) (mobyclient.ContainerRestartResult, error)
	ContainerRemove(ctx context.Context, containerID string, options mobyclient.ContainerRemoveOptions) (mobyclient.ContainerRemoveResult, error)
	ContainerList(ctx context.Context, options mobyclient.ContainerListOptions) (mobyclient.ContainerListResult, error)
	ContainerInspect(ctx context.Context, containerID string, options mobyclient.ContainerInspectOptions) (mobyclient.ContainerInspectResult, error)
	ContainerLogs(ctx context.Context, containerID string, options mobyclient.ContainerLogsOptions) (mobyclient.ContainerLogsResult, error)
	ExecCreate(ctx context.Context, containerID string, options mobyclient.ExecCreateOptions) (mobyclient.ExecCreateResult, error)
	ExecAttach(ctx context.Context, execID string, options mobyclient.ExecAttachOptions) (mobyclient.ExecAttachResult, error)
	ExecInspect(ctx context.Context, execID string, options mobyclient.ExecInspectOptions) (mobyclient.ExecInspectResult, error)
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
