package system

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (s *Service) Ping(ctx context.Context) (PingResult, error) {
	s.logger.Debug("pinging docker daemon")

	ping, err := s.providerClient.Ping(ctx)
	if err != nil {
		s.logger.Error("ping failed", zap.Error(err))
		return PingResult{}, fmt.Errorf("ping: %w", err)
	}

	return PingResult{
		APIVersion:     ping.APIVersion,
		OSType:         ping.OSType,
		Experimental:   ping.Experimental,
		BuilderVersion: ping.BuilderVersion,
	}, nil
}
