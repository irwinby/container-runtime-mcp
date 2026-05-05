package docker

import (
	"context"
	"fmt"
	"time"

	// Docker SDK is used for client construction because it resolves Docker
	// context and environment variables automatically. Moby client types are
	// used for the Engine API method signatures.
	dockerclient "github.com/docker/go-sdk/client"
	"github.com/irwinby/container-runtime-mcp/internal/provider/docker/container"
	"github.com/irwinby/container-runtime-mcp/internal/provider/docker/image"
	"github.com/irwinby/container-runtime-mcp/internal/provider/docker/system"
	"github.com/irwinby/container-runtime-mcp/internal/provider/docker/volume"
	mobyclient "github.com/moby/moby/client"
)

type dockerClient interface {
	ImagePull(ctx context.Context, refStr string, options mobyclient.ImagePullOptions) (mobyclient.ImagePullResponse, error)
	ImagePush(ctx context.Context, image string, options mobyclient.ImagePushOptions) (mobyclient.ImagePushResponse, error)
	ImageRemove(ctx context.Context, imageID string, options mobyclient.ImageRemoveOptions) (mobyclient.ImageRemoveResult, error)
	ImageTag(ctx context.Context, options mobyclient.ImageTagOptions) (mobyclient.ImageTagResult, error)
	ImageList(ctx context.Context, options mobyclient.ImageListOptions) (mobyclient.ImageListResult, error)
	ImageInspect(ctx context.Context, imageID string, inspectOpts ...mobyclient.ImageInspectOption) (mobyclient.ImageInspectResult, error)
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
	VolumeList(ctx context.Context, options mobyclient.VolumeListOptions) (mobyclient.VolumeListResult, error)
	VolumeInspect(ctx context.Context, volumeID string, options mobyclient.VolumeInspectOptions) (mobyclient.VolumeInspectResult, error)
	VolumeCreate(ctx context.Context, options mobyclient.VolumeCreateOptions) (mobyclient.VolumeCreateResult, error)
	VolumeRemove(ctx context.Context, volumeID string, options mobyclient.VolumeRemoveOptions) (mobyclient.VolumeRemoveResult, error)
	Ping(ctx context.Context, options mobyclient.PingOptions) (mobyclient.PingResult, error)
	Info(ctx context.Context, options mobyclient.InfoOptions) (mobyclient.SystemInfoResult, error)
	ServerVersion(ctx context.Context, options mobyclient.ServerVersionOptions) (mobyclient.ServerVersionResult, error)
	Close() error
}

// Compile-time check that the Docker SDK client satisfies the internal
// dockerClient interface, which is defined using Moby client types.
var _ dockerClient = (dockerclient.SDKClient)(nil)

type (
	ContainerProvider = *container.Provider
	ImageProvider     = *image.Provider
	VolumeProvider    = *volume.Provider
	SystemProvider    = *system.Provider
)

type Provider struct {
	ContainerProvider
	ImageProvider
	VolumeProvider
	SystemProvider

	client  dockerClient
	timeout time.Duration
}

func NewProvider(ctx context.Context, timeout time.Duration) (*Provider, error) {
	client, err := dockerclient.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("create docker client: %w", err)
	}

	return newProvider(client, timeout), nil
}

func newProvider(client dockerClient, timeout time.Duration) *Provider {
	withTimeout := func(ctx context.Context) (context.Context, context.CancelFunc) {
		if timeout > 0 {
			return context.WithTimeout(ctx, timeout)
		}

		return ctx, func() {}
	}

	return &Provider{
		ContainerProvider: container.NewProvider(client, withTimeout),
		ImageProvider:     image.NewProvider(client, withTimeout),
		VolumeProvider:    volume.NewProvider(client, withTimeout),
		SystemProvider:    system.NewProvider(client, withTimeout),
		client:            client,
		timeout:           timeout,
	}
}

func (p *Provider) Close() error {
	return p.client.Close()
}
