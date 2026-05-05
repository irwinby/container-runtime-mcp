package volume

import (
	"context"
	"fmt"
	"strings"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

// RemoveVolumeParams holds the parameters for removing a volume.
type RemoveVolumeParams struct {
	Name  string
	Force bool
}

func NewRemoveVolumeParams() RemoveVolumeParams {
	return RemoveVolumeParams{}
}

func (p RemoveVolumeParams) SetName(name string) RemoveVolumeParams {
	p.Name = strings.TrimSpace(name)
	return p
}

func (p RemoveVolumeParams) SetForce(force bool) RemoveVolumeParams {
	p.Force = force
	return p
}

// Validate checks that required fields are present and trims whitespace.
func (p RemoveVolumeParams) Validate() (RemoveVolumeParams, error) {
	if p.Name == "" {
		return RemoveVolumeParams{}, fmt.Errorf("name is required")
	}

	return p, nil
}

func (s *Service) RemoveVolume(ctx context.Context, params RemoveVolumeParams) error {
	err := s.policy.IsWriteAllowed()
	if err != nil {
		s.logger.Warn("remove volume blocked by policy", zap.Error(err))
		return fmt.Errorf("check if write is allowed: %w", err)
	}

	params, err = params.Validate()
	if err != nil {
		s.logger.Warn("remove volume validation failed", zap.Error(err))
		return fmt.Errorf("validate remove volume params: %w", err)
	}

	s.logger.Info("removing volume", zap.String("name", params.Name))

	err = s.providerClient.RemoveVolume(ctx, providers.RemoveVolumeParams{
		Name:  params.Name,
		Force: params.Force,
	})
	if err != nil {
		s.logger.Error("remove volume failed", zap.String("name", params.Name), zap.Error(err))
		return fmt.Errorf("remove volume: %w", err)
	}

	s.logger.Info("volume removed", zap.String("name", params.Name))

	return nil
}
