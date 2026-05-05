package image

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// InspectImageParams holds the parameters for inspecting an image.
type InspectImageParams struct {
	Ref string
}

func NewInspectImageParams() InspectImageParams {
	return InspectImageParams{}
}

func (p InspectImageParams) SetRef(ref string) InspectImageParams {
	p.Ref = strings.TrimSpace(ref)
	return p
}

// Validate checks that required fields are present and trims whitespace.
func (p InspectImageParams) Validate() (InspectImageParams, error) {
	if p.Ref == "" {
		return InspectImageParams{}, fmt.Errorf("ref is required")
	}

	return p, nil
}

func (s *Service) InspectImage(ctx context.Context, params InspectImageParams) (ImageInspect, error) {
	params, err := params.Validate()
	if err != nil {
		s.logger.Warn("inspect image validation failed", zap.Error(err))
		return ImageInspect{}, fmt.Errorf("validate inspect image params: %w", err)
	}

	s.logger.Debug("inspecting image", zap.String("ref", params.Ref))

	info, err := s.providerClient.InspectImage(ctx, providers.InspectImageParams{
		Ref: params.Ref,
	})
	if err != nil {
		s.logger.Error("inspect image failed", zap.String("ref", params.Ref), zap.Error(err))
		return ImageInspect{}, fmt.Errorf("inspect image: %w", err)
	}

	return ImageInspect{
		ID:           info.ID,
		RepoTags:     info.RepoTags,
		Size:         info.Size,
		Created:      info.Created,
		Architecture: info.Architecture,
		OS:           info.OS,
	}, nil
}
