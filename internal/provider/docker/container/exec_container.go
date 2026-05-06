package container

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

func (p *Provider) ExecContainer(ctx context.Context, params providers.ExecContainerParams) (providers.ExecContainerResult, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()

	createResult, err := p.client.ExecCreate(ctx, params.Name, client.ExecCreateOptions{
		User:         params.User,
		Privileged:   params.Privileged,
		TTY:          params.TTY,
		AttachStdin:  params.AttachStdin,
		AttachStdout: params.AttachStdout,
		AttachStderr: params.AttachStderr,
		Env:          params.Env,
		WorkingDir:   params.WorkingDir,
		Cmd:          params.Cmd,
	})
	if err != nil {
		return providers.ExecContainerResult{}, fmt.Errorf("create docker exec: %w", err)
	}

	attachResult, err := p.client.ExecAttach(ctx, createResult.ID, client.ExecAttachOptions{
		TTY: params.TTY,
	})
	if err != nil {
		return providers.ExecContainerResult{}, fmt.Errorf("attach docker exec: %w", err)
	}

	defer attachResult.Close()

	var stdoutBuffer, stderrBuffer bytes.Buffer
	var readErr error

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		if params.TTY {
			_, readErr = stdoutBuffer.ReadFrom(attachResult.Reader)
		} else {
			_, readErr = stdcopy.StdCopy(&stdoutBuffer, &stderrBuffer, attachResult.Reader)
		}
	}()

	if params.AttachStdin {
		_, err = io.WriteString(attachResult.Conn, params.Stdin)
		if err != nil {
			return providers.ExecContainerResult{}, fmt.Errorf("write exec stdin: %w", err)
		}

		err = attachResult.CloseWrite()
		if err != nil {
			return providers.ExecContainerResult{}, fmt.Errorf("close exec stdin: %w", err)
		}
	}

	wg.Wait()

	if readErr != nil {
		if params.TTY {
			return providers.ExecContainerResult{}, fmt.Errorf("read exec tty output: %w", readErr)
		}
		return providers.ExecContainerResult{}, fmt.Errorf("read exec output: %w", readErr)
	}

	inspectResult, err := p.client.ExecInspect(ctx, createResult.ID, client.ExecInspectOptions{})
	if err != nil {
		return providers.ExecContainerResult{}, fmt.Errorf("inspect docker exec: %w", err)
	}

	return providers.ExecContainerResult{
		ExecID:   createResult.ID,
		ExitCode: inspectResult.ExitCode,
		Stdout:   stdoutBuffer.String(),
		Stderr:   stderrBuffer.String(),
	}, nil
}
