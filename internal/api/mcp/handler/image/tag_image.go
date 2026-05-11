package image

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/image"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TagImageInput struct {
	Source string `json:"source" jsonschema:"source image reference"`
	Target string `json:"target" jsonschema:"target image reference"`
}

func (i *TagImageInput) Validate() error {
	if i == nil {
		return fmt.Errorf("tag image input is nil")
	}

	if i.Source == "" {
		return fmt.Errorf("source is required")
	}

	if i.Target == "" {
		return fmt.Errorf("target is required")
	}

	return nil
}

type TagImageOutput struct{}

func (h *ToolsHandler) TagImage(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input TagImageInput,
) (*mcp.CallToolResult, TagImageOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, TagImageOutput{}, fmt.Errorf("validate tag image input: %w", err)
	}

	params := image.NewTagImageParams().
		SetSource(input.Source).
		SetTarget(input.Target)

	err = h.imageService.TagImage(ctx, params)
	if err != nil {
		return nil, TagImageOutput{}, fmt.Errorf("tag image: %w", err)
	}

	return nil, TagImageOutput{}, nil
}
