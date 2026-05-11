package container

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ContainerLogsInput struct {
	Name       string `json:"name" jsonschema:"container name or ID"`
	Stdout     *bool  `json:"stdout,omitempty" jsonschema:"include stdout output"`
	Stderr     *bool  `json:"stderr,omitempty" jsonschema:"include stderr output"`
	Since      string `json:"since,omitempty" jsonschema:"show logs since timestamp or relative duration, e.g. 2024-01-01T00:00:00Z or 10m"`
	Timestamps bool   `json:"timestamps,omitempty" jsonschema:"include timestamps in output"`
	Tail       string `json:"tail,omitempty" jsonschema:"number of lines to show from the end of the logs, e.g. 100 or all"`
}

func (i *ContainerLogsInput) Validate() error {
	if i == nil {
		return fmt.Errorf("container logs input is nil")
	}

	if i.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

type ContainerLogsOutput struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func (h *ToolsHandler) ContainerLogs(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input ContainerLogsInput,
) (*mcp.CallToolResult, ContainerLogsOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, ContainerLogsOutput{}, fmt.Errorf("validate container logs input: %w", err)
	}

	params := container.NewContainerLogsParams().
		SetName(input.Name).
		SetSince(input.Since).
		SetTimestamps(input.Timestamps).
		SetTail(input.Tail)

	if input.Stdout != nil {
		params = params.SetStdout(*input.Stdout)
	}

	if input.Stderr != nil {
		params = params.SetStderr(*input.Stderr)
	}

	result, err := h.containerService.ContainerLogs(ctx, params)
	if err != nil {
		return nil, ContainerLogsOutput{}, fmt.Errorf("container logs: %w", err)
	}

	return nil, ContainerLogsOutput{
		Stdout: result.Stdout,
		Stderr: result.Stderr,
	}, nil
}
