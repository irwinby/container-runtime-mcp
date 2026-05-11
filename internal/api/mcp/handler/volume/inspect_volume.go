package volume

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/volume"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type InspectVolumeInput struct {
	Name string `json:"name" jsonschema:"volume name"`
}

func (i *InspectVolumeInput) Validate() error {
	if i == nil {
		return fmt.Errorf("inspect volume input is nil")
	}

	if i.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

type InspectVolumeOutput struct {
	Volume volume.VolumeInspect `json:"volume"`
}

func (h *ToolsHandler) InspectVolume(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input InspectVolumeInput,
) (*mcp.CallToolResult, InspectVolumeOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, InspectVolumeOutput{}, fmt.Errorf("validate inspect volume input: %w", err)
	}

	params := volume.NewInspectVolumeParams().
		SetName(input.Name)

	result, err := h.volumeService.InspectVolume(ctx, params)
	if err != nil {
		return nil, InspectVolumeOutput{}, fmt.Errorf("inspect volume: %w", err)
	}

	return nil, InspectVolumeOutput{Volume: result}, nil
}
