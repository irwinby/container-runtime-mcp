package volume

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/volume"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListVolumesInput struct {
	Dangling bool `json:"dangling,omitempty" jsonschema:"filter to show only dangling volumes"`
}

func (i *ListVolumesInput) Validate() error {
	if i == nil {
		return fmt.Errorf("list volumes input is nil")
	}

	return nil
}

type ListVolumesOutput struct {
	Volumes []volume.Volume `json:"volumes"`
}

func (h *ToolsHandler) ListVolumes(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input ListVolumesInput,
) (*mcp.CallToolResult, ListVolumesOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, ListVolumesOutput{}, fmt.Errorf("validate list volumes input: %w", err)
	}

	params := volume.NewListVolumesParams().
		SetDangling(input.Dangling)

	result, err := h.volumeService.ListVolumes(ctx, params)
	if err != nil {
		return nil, ListVolumesOutput{}, fmt.Errorf("list volumes: %w", err)
	}

	return nil, ListVolumesOutput{Volumes: result}, nil
}
