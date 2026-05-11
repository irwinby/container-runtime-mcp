package volume

import (
	"context"

	"github.com/irwinby/container-runtime-mcp/internal/service/volume"
	"github.com/irwinby/container-runtime-mcp/pkg/ptr"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type volumeService interface {
	CanWrite() bool
	ListVolumes(ctx context.Context, params volume.ListVolumesParams) ([]volume.Volume, error)
	InspectVolume(ctx context.Context, params volume.InspectVolumeParams) (volume.VolumeInspect, error)
	CreateVolume(ctx context.Context, params volume.CreateVolumeParams) (volume.VolumeInspect, error)
	RemoveVolume(ctx context.Context, params volume.RemoveVolumeParams) error
}

type ToolsHandler struct {
	volumeService volumeService
}

func NewToolsHandler(volumeService volumeService) *ToolsHandler {
	return &ToolsHandler{
		volumeService: volumeService,
	}
}

func (h *ToolsHandler) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_volumes",
		Title:       "List Volumes",
		Description: "List Docker volumes.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.ListVolumes)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "inspect_volume",
		Title:       "Inspect Volume",
		Description: "Inspect a Docker volume.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.InspectVolume)

	if !h.volumeService.CanWrite() {
		return
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_volume",
		Title:       "Create Volume",
		Description: "Create a Docker volume.",
		Annotations: &mcp.ToolAnnotations{
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.CreateVolume)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "remove_volume",
		Title:       "Remove Volume",
		Description: "Remove a Docker volume.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: ptr.Bool(true),
			OpenWorldHint:   ptr.Bool(false),
		},
	}, h.RemoveVolume)
}
