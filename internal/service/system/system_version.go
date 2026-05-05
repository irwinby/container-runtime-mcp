package system

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (s *Service) SystemVersion(ctx context.Context) (SystemVersion, error) {
	s.logger.Debug("fetching system version")

	version, err := s.providerClient.SystemVersion(ctx)
	if err != nil {
		s.logger.Error("system version failed", zap.Error(err))
		return SystemVersion{}, fmt.Errorf("system version: %w", err)
	}

	return SystemVersion{
		Version:       version.Version,
		APIVersion:    version.APIVersion,
		MinAPIVersion: version.MinAPIVersion,
		Os:            version.Os,
		Arch:          version.Arch,
		PlatformName:  version.PlatformName,
	}, nil
}
