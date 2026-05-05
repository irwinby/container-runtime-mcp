package image

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/image"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type InspectImageInput struct {
	Ref string `json:"ref" jsonschema:"image reference or ID"`
}

func (i *InspectImageInput) Validate() error {
	if i == nil {
		return fmt.Errorf("inspect image input is nil")
	}

	if i.Ref == "" {
		return fmt.Errorf("ref is required")
	}

	return nil
}

type InspectImageDetails struct {
	ID           string   `json:"id"`
	RepoTags     []string `json:"repo_tags"`
	Size         int64    `json:"size"`
	Created      string   `json:"created"`
	Architecture string   `json:"architecture"`
	OS           string   `json:"os"`
}

type InspectImageOutput struct {
	Image InspectImageDetails `json:"image"`
}

func (h *ToolsHandler) InspectImage(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input InspectImageInput,
) (*mcp.CallToolResult, InspectImageOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, InspectImageOutput{}, fmt.Errorf("validate inspect image input: %w", err)
	}

	params := image.NewInspectImageParams().
		SetRef(input.Ref)

	info, err := h.imageService.InspectImage(ctx, params)
	if err != nil {
		return nil, InspectImageOutput{}, fmt.Errorf("inspect image: %w", err)
	}

	return nil, InspectImageOutput{
		Image: InspectImageDetails{
			ID:           info.ID,
			RepoTags:     info.RepoTags,
			Size:         info.Size,
			Created:      info.Created,
			Architecture: info.Architecture,
			OS:           info.OS,
		},
	}, nil
}
