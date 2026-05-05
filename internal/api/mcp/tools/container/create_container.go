package container

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CreateContainerInput struct {
	Name  string `json:"name" jsonschema:"container name"`
	Image string `json:"image" jsonschema:"image to use for the container, for example nginx:latest"`
}

func (i *CreateContainerInput) Validate() error {
	if i == nil {
		return fmt.Errorf("create container input is nil")
	}

	if i.Name == "" {
		return fmt.Errorf("name is required")
	}

	if i.Image == "" {
		return fmt.Errorf("image is required")
	}

	return nil
}

type CreateContainerOutput struct {
	ID string `json:"id"`
}

func (h *ToolsHandler) CreateContainer(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input CreateContainerInput,
) (*mcp.CallToolResult, CreateContainerOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, CreateContainerOutput{}, fmt.Errorf("validate create container input: %w", err)
	}

	params := container.NewCreateContainerParams().
		SetName(input.Name).
		SetImage(input.Image)

	id, err := h.containerService.CreateContainer(ctx, params)
	if err != nil {
		return nil, CreateContainerOutput{}, fmt.Errorf("create container: %w", err)
	}

	return nil, CreateContainerOutput{ID: id}, nil
}
