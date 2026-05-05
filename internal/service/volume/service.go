package volume

import (
	"context"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	services "github.com/irwinby/container-runtime-mcp/internal/service"
	"go.uber.org/zap"
)

type providerClient interface {
	ListVolumes(ctx context.Context, params providers.ListVolumesParams) ([]providers.Volume, error)
	InspectVolume(ctx context.Context, params providers.InspectVolumeParams) (providers.VolumeInspect, error)
	CreateVolume(ctx context.Context, params providers.CreateVolumeParams) (providers.VolumeInspect, error)
	RemoveVolume(ctx context.Context, params providers.RemoveVolumeParams) error
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
