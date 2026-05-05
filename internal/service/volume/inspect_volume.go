package volume

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// InspectVolumeParams holds the parameters for inspecting a volume.
type InspectVolumeParams struct {
	Name string
}

func NewInspectVolumeParams() InspectVolumeParams {
	return InspectVolumeParams{}
}

func (p InspectVolumeParams) SetName(name string) InspectVolumeParams {
	p.Name = strings.TrimSpace(name)
	return p
}

// Validate checks that required fields are present and trims whitespace.
func (p InspectVolumeParams) Validate() (InspectVolumeParams, error) {
	if p.Name == "" {
		return InspectVolumeParams{}, fmt.Errorf("name is required")
	}

	return p, nil
}

type VolumeInspect struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	Labels     map[string]string `json:"labels"`
	Scope      string            `json:"scope"`
}

func (s *Service) InspectVolume(ctx context.Context, params InspectVolumeParams) (VolumeInspect, error) {
	params, err := params.Validate()
	if err != nil {
		s.logger.Warn("inspect volume validation failed", zap.Error(err))
		return VolumeInspect{}, fmt.Errorf("validate inspect volume params: %w", err)
	}

	s.logger.Debug("inspecting volume", zap.String("name", params.Name))

	result, err := s.providerClient.InspectVolume(ctx, providers.InspectVolumeParams{
		Name: params.Name,
	})
	if err != nil {
		s.logger.Error("inspect volume failed", zap.String("name", params.Name), zap.Error(err))
		return VolumeInspect{}, fmt.Errorf("inspect volume: %w", err)
	}

	return VolumeInspect{
		Name:       result.Name,
		Driver:     result.Driver,
		Mountpoint: result.Mountpoint,
		Labels:     result.Labels,
		Scope:      result.Scope,
	}, nil
}
