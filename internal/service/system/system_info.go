package system

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (s *Service) SystemInfo(ctx context.Context) (SystemInfo, error) {
	s.logger.Debug("fetching system info")

	info, err := s.providerClient.SystemInfo(ctx)
	if err != nil {
		s.logger.Error("system info failed", zap.Error(err))
		return SystemInfo{}, fmt.Errorf("system info: %w", err)
	}

	return SystemInfo{
		ID:                info.ID,
		Containers:        info.Containers,
		ContainersRunning: info.ContainersRunning,
		ContainersPaused:  info.ContainersPaused,
		ContainersStopped: info.ContainersStopped,
		Images:            info.Images,
		Driver:            info.Driver,
		DriverStatus:      info.DriverStatus,
		KernelVersion:     info.KernelVersion,
		OperatingSystem:   info.OperatingSystem,
		OSType:            info.OSType,
		Architecture:      info.Architecture,
		NCPU:              info.NCPU,
		MemTotal:          info.MemTotal,
		ServerVersion:     info.ServerVersion,
	}, nil
}
