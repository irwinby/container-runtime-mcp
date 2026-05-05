package container

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type InspectContainerInput struct {
	Name string `json:"name" jsonschema:"container name or ID"`
}

func (i *InspectContainerInput) Validate() error {
	if i == nil {
		return fmt.Errorf("inspect container input is nil")
	}

	if i.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

type InspectContainerDetails struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	State        string   `json:"state"`
	Status       string   `json:"status"`
	Created      string   `json:"created"`
	Path         string   `json:"path"`
	Args         []string `json:"args"`
	RestartCount int      `json:"restart_count"`
}

type InspectContainerOutput struct {
	Container InspectContainerDetails `json:"container"`
}

func (h *ToolsHandler) InspectContainer(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input InspectContainerInput,
) (*mcp.CallToolResult, InspectContainerOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, InspectContainerOutput{}, fmt.Errorf("validate inspect container input: %w", err)
	}

	params := container.NewInspectContainerParams().
		SetName(input.Name)

	info, err := h.containerService.InspectContainer(ctx, params)
	if err != nil {
		return nil, InspectContainerOutput{}, fmt.Errorf("inspect container: %w", err)
	}

	return nil, InspectContainerOutput{
		Container: InspectContainerDetails{
			ID:           info.ID,
			Name:         info.Name,
			Image:        info.Image,
			State:        info.State,
			Status:       info.Status,
			Created:      info.Created,
			Path:         info.Path,
			Args:         info.Args,
			RestartCount: info.RestartCount,
		},
	}, nil
}
