package image

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/image"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListImagesInput struct {
	All        bool `json:"all,omitempty" jsonschema:"list all images, including intermediate layers"`
	SharedSize bool `json:"shared_size,omitempty" jsonschema:"include shared size computation"`
}

type ListImagesItem struct {
	ID         string   `json:"id"`
	RepoTags   []string `json:"repo_tags"`
	Size       int64    `json:"size"`
	Created    int64    `json:"created"`
	Containers int64    `json:"containers"`
}

type ListImagesOutput struct {
	Images []ListImagesItem `json:"images"`
}

func (h *ToolsHandler) ListImages(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input ListImagesInput,
) (*mcp.CallToolResult, ListImagesOutput, error) {
	params := image.NewListImagesParams().
		SetAll(input.All).
		SetSharedSize(input.SharedSize)

	images, err := h.imageService.ListImages(ctx, params)
	if err != nil {
		return nil, ListImagesOutput{}, fmt.Errorf("list images: %w", err)
	}

	output := ListImagesOutput{
		Images: make([]ListImagesItem, 0, len(images)),
	}

	for _, image := range images {
		output.Images = append(output.Images, ListImagesItem{
			ID:         image.ID,
			RepoTags:   image.RepoTags,
			Size:       image.Size,
			Created:    image.Created,
			Containers: image.Containers,
		})
	}

	return nil, output, nil
}
