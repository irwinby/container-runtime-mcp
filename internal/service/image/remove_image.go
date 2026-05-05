package image

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"go.uber.org/zap"
)

// RemoveImageParams holds the parameters for removing an image.
type RemoveImageParams struct {
	Ref           string
	Force         bool
	PruneChildren bool
	Platform      *ocispec.Platform
}

func NewRemoveImageParams() RemoveImageParams {
	return RemoveImageParams{}
}

func (p RemoveImageParams) SetRef(ref string) RemoveImageParams {
	p.Ref = strings.TrimSpace(ref)
	return p
}

func (p RemoveImageParams) SetForce(force bool) RemoveImageParams {
	p.Force = force
	return p
}

func (p RemoveImageParams) SetPruneChildren(pruneChildren bool) RemoveImageParams {
	p.PruneChildren = pruneChildren
	return p
}

func (p RemoveImageParams) SetPlatform(platform *ocispec.Platform) RemoveImageParams {
	p.Platform = platform
	return p
}

// Validate checks that required fields are present and trims whitespace.
func (p RemoveImageParams) Validate() (RemoveImageParams, error) {
	if p.Ref == "" {
		return RemoveImageParams{}, fmt.Errorf("ref is required")
	}

	return p, nil
}

func (s *Service) RemoveImage(ctx context.Context, params RemoveImageParams) error {
	err := s.policy.IsWriteAllowed()
	if err != nil {
		s.logger.Warn("remove image blocked by policy", zap.Error(err))
		return fmt.Errorf("check if write is allowed: %w", err)
	}

	params, err = params.Validate()
	if err != nil {
		s.logger.Warn("remove image validation failed", zap.Error(err))
		return fmt.Errorf("validate remove image params: %w", err)
	}

	s.logger.Info("removing image", zap.String("ref", params.Ref))

	err = s.providerClient.RemoveImage(ctx, providers.RemoveImageParams{
		Ref:           params.Ref,
		Force:         params.Force,
		PruneChildren: params.PruneChildren,
		Platform:      params.Platform,
	})
	if err != nil {
		s.logger.Error("remove image failed", zap.String("ref", params.Ref), zap.Error(err))
		return fmt.Errorf("remove image: %w", err)
	}

	s.logger.Info("image removed", zap.String("ref", params.Ref))

	return nil
}
