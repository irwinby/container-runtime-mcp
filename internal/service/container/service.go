package container

import (
	"context"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	services "github.com/irwinby/container-runtime-mcp/internal/service"
	"go.uber.org/zap"
)

type providerClient interface {
	CreateContainer(ctx context.Context, params providers.CreateContainerParams) (string, error)
	RemoveContainer(ctx context.Context, params providers.RemoveContainerParams) error
	ListContainers(ctx context.Context, params providers.ListContainersParams) ([]providers.Container, error)
	InspectContainer(ctx context.Context, params providers.InspectContainerParams) (providers.ContainerInspect, error)
	StartContainer(ctx context.Context, params providers.StartContainerParams) error
	StopContainer(ctx context.Context, params providers.StopContainerParams) error
	RestartContainer(ctx context.Context, params providers.RestartContainerParams) error
	ContainerLogs(ctx context.Context, params providers.ContainerLogsParams) (providers.ContainerLogsResult, error)
	ExecContainer(ctx context.Context, params providers.ExecContainerParams) (providers.ExecContainerResult, error)
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
