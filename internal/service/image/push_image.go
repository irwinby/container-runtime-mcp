package image

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"go.uber.org/zap"
)

// PushImageParams holds the parameters for pushing an image.
type PushImageParams struct {
	Ref      string
	All      bool
	Platform *ocispec.Platform
}

func NewPushImageParams() PushImageParams {
	return PushImageParams{}
}

func (p PushImageParams) SetRef(ref string) PushImageParams {
	p.Ref = strings.TrimSpace(ref)

	return p
}

func (p PushImageParams) SetAll(all bool) PushImageParams {
	p.All = all

	return p
}

func (p PushImageParams) SetPlatform(platform *ocispec.Platform) PushImageParams {
	p.Platform = platform

	return p
}

// Validate checks that required fields are present, trims whitespace, and parses the platform.
func (p PushImageParams) Validate() (PushImageParams, error) {
	if p.Ref == "" {
		return PushImageParams{}, fmt.Errorf("ref is required")
	}

	return p, nil
}

func (s *Service) PushImage(ctx context.Context, params PushImageParams) error {
	err := s.policy.IsWriteAllowed()
	if err != nil {
		s.logger.Warn("push image blocked by policy", zap.Error(err))
		return fmt.Errorf("check if write is allowed: %w", err)
	}

	params, err = params.Validate()
	if err != nil {
		s.logger.Warn("push image validation failed", zap.Error(err))
		return fmt.Errorf("validate push image params: %w", err)
	}

	s.logger.Info("pushing image", zap.String("ref", params.Ref))

	err = s.providerClient.PushImage(ctx, providers.PushImageParams{
		Ref:      params.Ref,
		All:      params.All,
		Platform: params.Platform,
	})
	if err != nil {
		s.logger.Error("push image failed", zap.String("ref", params.Ref), zap.Error(err))
		return fmt.Errorf("push image: %w", err)
	}

	s.logger.Info("image pushed", zap.String("ref", params.Ref))

	return nil
}
