package image

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"go.uber.org/zap"
)

// PullImageParams holds the parameters for pulling an image.
type PullImageParams struct {
	Ref      string
	All      bool
	Platform *ocispec.Platform
}

func NewPullImageParams() PullImageParams {
	return PullImageParams{}
}

func (p PullImageParams) SetRef(ref string) PullImageParams {
	p.Ref = strings.TrimSpace(ref)

	return p
}

func (p PullImageParams) SetAll(all bool) PullImageParams {
	p.All = all

	return p
}

func (p PullImageParams) SetPlatform(platform *ocispec.Platform) PullImageParams {
	p.Platform = platform

	return p
}

// Validate checks that required fields are present, trims whitespace, and parses the platform.
func (p PullImageParams) Validate() (PullImageParams, error) {
	if p.Ref == "" {
		return PullImageParams{}, fmt.Errorf("ref is required")
	}

	return p, nil
}

func (s *Service) PullImage(ctx context.Context, params PullImageParams) error {
	err := s.policy.IsWriteAllowed()
	if err != nil {
		s.logger.Warn("pull image blocked by policy", zap.Error(err))
		return fmt.Errorf("check if write is allowed: %w", err)
	}

	params, err = params.Validate()
	if err != nil {
		s.logger.Warn("pull image validation failed", zap.Error(err))
		return fmt.Errorf("validate pull image params: %w", err)
	}

	s.logger.Info("pulling image", zap.String("ref", params.Ref))

	err = s.providerClient.PullImage(ctx, providers.PullImageParams{
		Ref:      params.Ref,
		All:      params.All,
		Platform: params.Platform,
	})
	if err != nil {
		s.logger.Error("pull image failed", zap.String("ref", params.Ref), zap.Error(err))
		return fmt.Errorf("pull image: %w", err)
	}

	s.logger.Info("image pulled", zap.String("ref", params.Ref))

	return nil
}
