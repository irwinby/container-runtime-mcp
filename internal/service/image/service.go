package image

import (
	"context"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	services "github.com/irwinby/container-runtime-mcp/internal/service"
	"go.uber.org/zap"
)

type providerClient interface {
	PullImage(ctx context.Context, params providers.PullImageParams) error
	PushImage(ctx context.Context, params providers.PushImageParams) error
	ListImages(ctx context.Context, params providers.ListImagesParams) ([]providers.Image, error)
	InspectImage(ctx context.Context, params providers.InspectImageParams) (providers.ImageInspect, error)
	RemoveImage(ctx context.Context, params providers.RemoveImageParams) error
	TagImage(ctx context.Context, params providers.TagImageParams) error
}

type Service struct {
	providerClient providerClient
	policy         services.Policy
	logger         *zap.Logger
}

func NewService(providerClient providerClient, policy services.Policy, logger *zap.Logger) *Service {
	return &Service{
		providerClient: providerClient,
		policy:         policy,
		logger:         logger,
	}
}

// CanWrite reports whether write operations are allowed.
func (s *Service) CanWrite() bool {
	return !s.policy.ReadOnly
}
