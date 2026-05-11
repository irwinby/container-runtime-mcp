package image

import (
	"context"

	"github.com/irwinby/container-runtime-mcp/internal/service/image"
	"github.com/irwinby/container-runtime-mcp/pkg/ptr"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type imageService interface {
	CanWrite() bool
	PullImage(ctx context.Context, params image.PullImageParams) error
	PushImage(ctx context.Context, params image.PushImageParams) error
	ListImages(ctx context.Context, params image.ListImagesParams) ([]image.Image, error)
	InspectImage(ctx context.Context, params image.InspectImageParams) (image.ImageInspect, error)
	RemoveImage(ctx context.Context, params image.RemoveImageParams) error
	TagImage(ctx context.Context, params image.TagImageParams) error
}

type ToolsHandler struct {
	imageService imageService
}

func NewToolsHandler(imageService imageService) *ToolsHandler {
	return &ToolsHandler{
		imageService: imageService,
	}
}

func (h *ToolsHandler) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_images",
		Title:       "List Images",
		Description: "List Docker images.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.ListImages)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "inspect_image",
		Title:       "Inspect Image",
		Description: "Inspect a Docker image.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.InspectImage)

	if !h.imageService.CanWrite() {
		return
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "pull_image",
		Title:       "Pull Image",
		Description: "Pull a Docker image from a registry by reference.",
		Annotations: &mcp.ToolAnnotations{
			OpenWorldHint: ptr.Bool(true),
		},
	}, h.PullImage)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "push_image",
		Title:       "Push Image",
		Description: "Push a Docker image to a registry by reference.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: ptr.Bool(true),
			OpenWorldHint:   ptr.Bool(true),
		},
	}, h.PushImage)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "remove_image",
		Title:       "Remove Image",
		Description: "Remove a Docker image.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: ptr.Bool(true),
			OpenWorldHint:   ptr.Bool(false),
		},
	}, h.RemoveImage)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "tag_image",
		Title:       "Tag Image",
		Description: "Tag a Docker image.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: ptr.Bool(true),
			OpenWorldHint:   ptr.Bool(false),
		},
	}, h.TagImage)
}
