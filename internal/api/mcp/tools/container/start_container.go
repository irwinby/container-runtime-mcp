package container

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type StartContainerInput struct {
	Name string `json:"name" jsonschema:"container name or ID"`
}

func (i *StartContainerInput) Validate() error {
	if i == nil {
		return fmt.Errorf("start container input is nil")
	}

	if i.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

type StartContainerOutput struct{}

func (h *ToolsHandler) StartContainer(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input StartContainerInput,
) (*mcp.CallToolResult, StartContainerOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, StartContainerOutput{}, fmt.Errorf("validate start container input: %w", err)
	}

	params := container.NewStartContainerParams().
		SetName(input.Name)

	err = h.containerService.StartContainer(ctx, params)
	if err != nil {
		return nil, StartContainerOutput{}, fmt.Errorf("start container: %w", err)
	}

	return nil, StartContainerOutput{}, nil
}
