package volume

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/volume"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type RemoveVolumeInput struct {
	Name  string `json:"name" jsonschema:"volume name"`
	Force bool   `json:"force,omitempty" jsonschema:"force removal"`
}

func (i *RemoveVolumeInput) Validate() error {
	if i == nil {
		return fmt.Errorf("remove volume input is nil")
	}

	if i.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

type RemoveVolumeOutput struct{}

func (h *ToolsHandler) RemoveVolume(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input RemoveVolumeInput,
) (*mcp.CallToolResult, RemoveVolumeOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, RemoveVolumeOutput{}, fmt.Errorf("validate remove volume input: %w", err)
	}

	params := volume.NewRemoveVolumeParams().
		SetName(input.Name).
		SetForce(input.Force)

	err = h.volumeService.RemoveVolume(ctx, params)
	if err != nil {
		return nil, RemoveVolumeOutput{}, fmt.Errorf("remove volume: %w", err)
	}

	return nil, RemoveVolumeOutput{}, nil
}
