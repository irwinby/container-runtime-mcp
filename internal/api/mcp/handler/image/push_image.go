package image

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/image"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type PushImageInput struct {
	Ref      string            `json:"ref" jsonschema:"image reference to push, for example registry.example.com/app:latest"`
	All      bool              `json:"all,omitempty" jsonschema:"push all tags for the image repository"`
	Platform *ocispec.Platform `json:"platform,omitempty" jsonschema:"optional platform, for example linux/amd64"`
}

func (i *PushImageInput) Validate() error {
	if i == nil {
		return fmt.Errorf("push image input is nil")
	}

	if i.Ref == "" {
		return fmt.Errorf("ref is required")
	}

	return nil
}

type PushImageOutput struct{}

func (h *ToolsHandler) PushImage(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input PushImageInput,
) (*mcp.CallToolResult, PushImageOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, PushImageOutput{}, fmt.Errorf("validate push image input: %w", err)
	}

	params := image.NewPushImageParams().
		SetRef(input.Ref).
		SetAll(input.All).
		SetPlatform(input.Platform)

	err = h.imageService.PushImage(ctx, params)
	if err != nil {
		return nil, PushImageOutput{}, fmt.Errorf("push image: %w", err)
	}

	return nil, PushImageOutput{}, nil
}
