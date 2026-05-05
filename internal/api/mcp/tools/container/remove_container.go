package container

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type RemoveContainerInput struct {
	Name          string `json:"name" jsonschema:"container name or ID"`
	Force         bool   `json:"force,omitempty" jsonschema:"force removal of a running container"`
	RemoveVolumes bool   `json:"remove_volumes,omitempty" jsonschema:"remove anonymous volumes associated with the container"`
	RemoveLinks   bool   `json:"remove_links,omitempty" jsonschema:"remove the specified link"`
}

func (i *RemoveContainerInput) Validate() error {
	if i == nil {
		return nil
	}

	if i.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

type RemoveContainerOutput struct{}

func (h *ToolsHandler) RemoveContainer(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input RemoveContainerInput,
) (*mcp.CallToolResult, RemoveContainerOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, RemoveContainerOutput{}, fmt.Errorf("validate remove container input: %w", err)
	}

	params := container.NewRemoveContainerParams().
		SetName(input.Name).
		SetForce(input.Force).
		SetRemoveVolumes(input.RemoveVolumes).
		SetRemoveLinks(input.RemoveLinks)

	err = h.containerService.RemoveContainer(ctx, params)
	if err != nil {
		return nil, RemoveContainerOutput{}, fmt.Errorf("remove container: %w", err)
	}

	return nil, RemoveContainerOutput{}, nil
}
