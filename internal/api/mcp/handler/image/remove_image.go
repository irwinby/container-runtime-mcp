package image

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/image"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type RemoveImageInput struct {
	Ref           string            `json:"ref" jsonschema:"image reference or ID"`
	Force         bool              `json:"force,omitempty" jsonschema:"force removal"`
	PruneChildren bool              `json:"prune_children,omitempty" jsonschema:"remove untagged parent images"`
	Platform      *ocispec.Platform `json:"platform,omitempty" jsonschema:"optional platform, for example linux/amd64"`
}

func (i *RemoveImageInput) Validate() error {
	if i == nil {
		return fmt.Errorf("remove image input is nil")
	}

	if i.Ref == "" {
		return fmt.Errorf("ref is required")
	}

	return nil
}

type RemoveImageOutput struct{}

func (h *ToolsHandler) RemoveImage(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input RemoveImageInput,
) (*mcp.CallToolResult, RemoveImageOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, RemoveImageOutput{}, fmt.Errorf("validate remove image input: %w", err)
	}

	params := image.NewRemoveImageParams().
		SetRef(input.Ref).
		SetForce(input.Force).
		SetPruneChildren(input.PruneChildren).
		SetPlatform(input.Platform)

	err = h.imageService.RemoveImage(ctx, params)
	if err != nil {
		return nil, RemoveImageOutput{}, fmt.Errorf("remove image: %w", err)
	}

	return nil, RemoveImageOutput{}, nil
}
