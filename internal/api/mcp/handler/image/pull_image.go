package image

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/image"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type PullImageInput struct {
	Ref      string            `json:"ref" jsonschema:"image reference to pull, for example nginx:latest"`
	All      bool              `json:"all,omitempty" jsonschema:"pull all tagged images in the repository"`
	Platform *ocispec.Platform `json:"platform,omitempty" jsonschema:"optional platform, for example linux/amd64"`
}

func (i *PullImageInput) Validate() error {
	if i == nil {
		return fmt.Errorf("pull image input is nil")
	}

	if i.Ref == "" {
		return fmt.Errorf("ref is required")
	}

	return nil
}

type PullImageOutput struct{}

func (h *ToolsHandler) PullImage(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input PullImageInput,
) (*mcp.CallToolResult, PullImageOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, PullImageOutput{}, fmt.Errorf("validate pull image input: %w", err)
	}

	params := image.NewPullImageParams().
		SetRef(input.Ref).
		SetAll(input.All).
		SetPlatform(input.Platform)

	err = h.imageService.PullImage(ctx, params)
	if err != nil {
		return nil, PullImageOutput{}, fmt.Errorf("pull image: %w", err)
	}

	return nil, PullImageOutput{}, nil
}
