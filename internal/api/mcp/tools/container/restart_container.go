package container

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type RestartContainerInput struct {
	Name           string `json:"name" jsonschema:"container name or ID"`
	Signal         string `json:"signal,omitempty" jsonschema:"signal to send, e.g. SIGTERM"`
	TimeoutSeconds *int   `json:"timeout_seconds,omitempty" jsonschema:"timeout in seconds before forceful kill; nil for default, -1 for indefinite, 0 for immediate"`
}

func (i *RestartContainerInput) Validate() error {
	if i == nil {
		return fmt.Errorf("restart container input is nil")
	}

	if i.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

type RestartContainerOutput struct{}

func (h *ToolsHandler) RestartContainer(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input RestartContainerInput,
) (*mcp.CallToolResult, RestartContainerOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, RestartContainerOutput{}, fmt.Errorf("validate restart container input: %w", err)
	}

	params := container.NewRestartContainerParams().
		SetName(input.Name).
		SetSignal(input.Signal).
		SetTimeoutSeconds(input.TimeoutSeconds)

	err = h.containerService.RestartContainer(ctx, params)
	if err != nil {
		return nil, RestartContainerOutput{}, fmt.Errorf("restart container: %w", err)
	}

	return nil, RestartContainerOutput{}, nil
}
