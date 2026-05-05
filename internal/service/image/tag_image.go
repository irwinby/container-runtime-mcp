package image

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// TagImageParams holds the parameters for tagging an image.
type TagImageParams struct {
	Source string
	Target string
}

func NewTagImageParams() TagImageParams {
	return TagImageParams{}
}

func (p TagImageParams) SetSource(source string) TagImageParams {
	p.Source = strings.TrimSpace(source)
	return p
}

func (p TagImageParams) SetTarget(target string) TagImageParams {
	p.Target = strings.TrimSpace(target)
	return p
}

// Validate checks that required fields are present and trims whitespace.
func (p TagImageParams) Validate() (TagImageParams, error) {
	if p.Source == "" {
		return TagImageParams{}, fmt.Errorf("source is required")
	}

	if p.Target == "" {
		return TagImageParams{}, fmt.Errorf("target is required")
	}

	return p, nil
}

func (s *Service) TagImage(ctx context.Context, params TagImageParams) error {
	err := s.policy.IsWriteAllowed()
	if err != nil {
		s.logger.Warn("tag image blocked by policy", zap.Error(err))
		return fmt.Errorf("check if write is allowed: %w", err)
	}

	params, err = params.Validate()
	if err != nil {
		s.logger.Warn("tag image validation failed", zap.Error(err))
		return fmt.Errorf("validate tag image params: %w", err)
	}

	s.logger.Info("tagging image", zap.String("source", params.Source), zap.String("target", params.Target))

	err = s.providerClient.TagImage(ctx, providers.TagImageParams{
		Source: params.Source,
		Target: params.Target,
	})
	if err != nil {
		s.logger.Error("tag image failed", zap.String("source", params.Source), zap.Error(err))
		return fmt.Errorf("tag image: %w", err)
	}

	s.logger.Info("image tagged", zap.String("target", params.Target))

	return nil
}
