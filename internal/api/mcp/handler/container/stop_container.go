package container

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type StopContainerInput struct {
	Name           string `json:"name" jsonschema:"container name or ID"`
	Signal         string `json:"signal,omitempty" jsonschema:"signal to send, e.g. SIGTERM"`
	TimeoutSeconds *int   `json:"timeout_seconds,omitempty" jsonschema:"timeout in seconds before forceful kill; nil for default, -1 for indefinite, 0 for immediate"`
}

func (i *StopContainerInput) Validate() error {
	if i == nil {
		return fmt.Errorf("stop container input is nil")
	}

	if i.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

type StopContainerOutput struct{}

func (h *ToolsHandler) StopContainer(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input StopContainerInput,
) (*mcp.CallToolResult, StopContainerOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, StopContainerOutput{}, fmt.Errorf("validate stop container input: %w", err)
	}

	params := container.NewStopContainerParams().
		SetName(input.Name).
		SetSignal(input.Signal).
		SetTimeoutSeconds(input.TimeoutSeconds)

	err = h.containerService.StopContainer(ctx, params)
	if err != nil {
		return nil, StopContainerOutput{}, fmt.Errorf("stop container: %w", err)
	}

	return nil, StopContainerOutput{}, nil
}
