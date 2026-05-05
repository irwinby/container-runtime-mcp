package volume

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/volume"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CreateVolumeInput struct {
	Name       string            `json:"name,omitempty" jsonschema:"volume name"`
	Driver     string            `json:"driver,omitempty" jsonschema:"volume driver"`
	DriverOpts map[string]string `json:"driver_opts,omitempty" jsonschema:"driver specific options"`
	Labels     map[string]string `json:"labels,omitempty" jsonschema:"volume labels"`
}

func (i *CreateVolumeInput) Validate() error {
	if i == nil {
		return fmt.Errorf("create volume input is nil")
	}

	return nil
}

type CreateVolumeOutput struct {
	Volume volume.VolumeInspect `json:"volume"`
}

func (h *ToolsHandler) CreateVolume(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input CreateVolumeInput,
) (*mcp.CallToolResult, CreateVolumeOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, CreateVolumeOutput{}, fmt.Errorf("validate create volume input: %w", err)
	}

	params := volume.NewCreateVolumeParams().
		SetName(input.Name).
		SetDriver(input.Driver).
		SetDriverOpts(input.DriverOpts).
		SetLabels(input.Labels)

	result, err := h.volumeService.CreateVolume(ctx, params)
	if err != nil {
		return nil, CreateVolumeOutput{}, fmt.Errorf("create volume: %w", err)
	}

	return nil, CreateVolumeOutput{Volume: result}, nil
}
