package system

import (
	"context"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"go.uber.org/zap"
)

type providerClient interface {
	SystemInfo(ctx context.Context) (providers.SystemInfo, error)
	SystemVersion(ctx context.Context) (providers.SystemVersion, error)
	Ping(ctx context.Context) (providers.PingResult, error)
}

type Service struct {
	providerClient providerClient
	logger         *zap.Logger
}

func NewService(providerClient providerClient, logger *zap.Logger) *Service {
	return &Service{
		providerClient: providerClient,
		logger:         logger,
	}
}
