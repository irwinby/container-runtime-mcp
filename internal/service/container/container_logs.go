package container

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// ContainerLogsParams holds the parameters for retrieving container logs.
type ContainerLogsParams struct {
	Name       string
	Stdout     bool
	Stderr     bool
	Since      string
	Timestamps bool
	Tail       string
}

func NewContainerLogsParams() ContainerLogsParams {
	return ContainerLogsParams{
		Stdout: true,
		Stderr: true,
	}
}

func (p ContainerLogsParams) SetName(name string) ContainerLogsParams {
	p.Name = strings.TrimSpace(name)
	return p
}

func (p ContainerLogsParams) SetStdout(stdout bool) ContainerLogsParams {
	p.Stdout = stdout
	return p
}

func (p ContainerLogsParams) SetStderr(stderr bool) ContainerLogsParams {
	p.Stderr = stderr
	return p
}

func (p ContainerLogsParams) SetSince(since string) ContainerLogsParams {
	p.Since = strings.TrimSpace(since)
	return p
}

func (p ContainerLogsParams) SetTimestamps(timestamps bool) ContainerLogsParams {
	p.Timestamps = timestamps
	return p
}

func (p ContainerLogsParams) SetTail(tail string) ContainerLogsParams {
	p.Tail = strings.TrimSpace(tail)
	return p
}

// Validate checks that required fields are present and trims whitespace.
func (p ContainerLogsParams) Validate() (ContainerLogsParams, error) {
	if p.Name == "" {
		return ContainerLogsParams{}, fmt.Errorf("name is required")
	}

	if !p.Stdout && !p.Stderr {
		return ContainerLogsParams{}, fmt.Errorf("at least one of stdout or stderr must be true")
	}

	return p, nil
}

func (s *Service) ContainerLogs(ctx context.Context, params ContainerLogsParams) (ContainerLogsResult, error) {
	params, err := params.Validate()
	if err != nil {
		s.logger.Warn("container logs validation failed", zap.Error(err))
		return ContainerLogsResult{}, fmt.Errorf("validate container logs params: %w", err)
	}

	s.logger.Debug("fetching container logs", zap.String("name", params.Name))

	result, err := s.providerClient.ContainerLogs(ctx, providers.ContainerLogsParams{
		Name:       params.Name,
		Stdout:     params.Stdout,
		Stderr:     params.Stderr,
		Since:      params.Since,
		Timestamps: params.Timestamps,
		Tail:       params.Tail,
	})
	if err != nil {
		s.logger.Error("container logs failed", zap.String("name", params.Name), zap.Error(err))
		return ContainerLogsResult{}, fmt.Errorf("container logs: %w", err)
	}

	return ContainerLogsResult{
		Stdout: result.Stdout,
		Stderr: result.Stderr,
	}, nil
}
