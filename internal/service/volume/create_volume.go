package volume

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// CreateVolumeParams holds the parameters for creating a volume.
type CreateVolumeParams struct {
	Name       string
	Driver     string
	DriverOpts map[string]string
	Labels     map[string]string
}

func NewCreateVolumeParams() CreateVolumeParams {
	return CreateVolumeParams{}
}

func (p CreateVolumeParams) SetName(name string) CreateVolumeParams {
	p.Name = strings.TrimSpace(name)
	return p
}

func (p CreateVolumeParams) SetDriver(driver string) CreateVolumeParams {
	p.Driver = strings.TrimSpace(driver)
	return p
}

func (p CreateVolumeParams) SetDriverOpts(driverOpts map[string]string) CreateVolumeParams {
	p.DriverOpts = driverOpts
	return p
}

func (p CreateVolumeParams) SetLabels(labels map[string]string) CreateVolumeParams {
	p.Labels = labels
	return p
}

func (s *Service) CreateVolume(ctx context.Context, params CreateVolumeParams) (VolumeInspect, error) {
	err := s.policy.IsWriteAllowed()
	if err != nil {
		s.logger.Warn("create volume blocked by policy", zap.Error(err))
		return VolumeInspect{}, fmt.Errorf("check if write is allowed: %w", err)
	}

	s.logger.Info("creating volume", zap.String("name", params.Name))

	result, err := s.providerClient.CreateVolume(ctx, providers.CreateVolumeParams{
		Name:       params.Name,
		Driver:     params.Driver,
		DriverOpts: params.DriverOpts,
		Labels:     params.Labels,
	})
	if err != nil {
		s.logger.Error("create volume failed", zap.String("name", params.Name), zap.Error(err))
		return VolumeInspect{}, fmt.Errorf("create volume: %w", err)
	}

	s.logger.Info("volume created", zap.String("name", result.Name))

	return VolumeInspect{
		Name:       result.Name,
		Driver:     result.Driver,
		Mountpoint: result.Mountpoint,
		Labels:     result.Labels,
		Scope:      result.Scope,
	}, nil
}
