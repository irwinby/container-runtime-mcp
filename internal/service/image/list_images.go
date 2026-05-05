package image

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// ListImagesParams holds the parameters for listing images.
type ListImagesParams struct {
	All        bool
	SharedSize bool
}

func NewListImagesParams() ListImagesParams {
	return ListImagesParams{}
}

func (p ListImagesParams) SetAll(all bool) ListImagesParams {
	p.All = all
	return p
}

func (p ListImagesParams) SetSharedSize(sharedSize bool) ListImagesParams {
	p.SharedSize = sharedSize
	return p
}

func (s *Service) ListImages(ctx context.Context, params ListImagesParams) ([]Image, error) {
	s.logger.Debug("listing images")

	images, err := s.providerClient.ListImages(ctx, providers.ListImagesParams{
		All:        params.All,
		SharedSize: params.SharedSize,
	})
	if err != nil {
		s.logger.Error("list images failed", zap.Error(err))
		return nil, fmt.Errorf("list images: %w", err)
	}

	result := make([]Image, 0, len(images))

	for _, image := range images {
		result = append(result, Image{
			ID:         image.ID,
			RepoTags:   image.RepoTags,
			Size:       image.Size,
			Created:    image.Created,
			Containers: image.Containers,
		})
	}

	return result, nil
}
