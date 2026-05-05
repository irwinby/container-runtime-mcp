package container

import (
	"bytes"
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

func (p *Provider) ContainerLogs(ctx context.Context, params providers.ContainerLogsParams) (providers.ContainerLogsResult, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()

	result, err := p.client.ContainerLogs(ctx, params.Name, client.ContainerLogsOptions{
		ShowStdout: params.Stdout,
		ShowStderr: params.Stderr,
		Since:      params.Since,
		Timestamps: params.Timestamps,
		Tail:       params.Tail,
		Follow:     false,
	})
	if err != nil {
		return providers.ContainerLogsResult{}, fmt.Errorf("get docker container logs: %w", err)
	}

	defer result.Close()

	var stdoutBuf, stderrBuf bytes.Buffer

	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, result)
	if err != nil {
		return providers.ContainerLogsResult{}, fmt.Errorf("read container logs: %w", err)
	}

	return providers.NewContainerLogsResult().
		SetStdout(stdoutBuf.String()).
		SetStderr(stderrBuf.String()), nil
}
