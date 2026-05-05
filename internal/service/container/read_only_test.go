package container

import (
	"context"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	services "github.com/irwinby/container-runtime-mcp/internal/service"
	containermock "github.com/irwinby/container-runtime-mcp/internal/service/container/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestService_ReadOnly(t *testing.T) {
	policy := services.NewPolicy(true)

	mockClient := containermock.NewMockProviderClient(t)
	service := NewService(mockClient, policy, zap.NewNop())

	ctx := context.Background()

	t.Run("CreateContainer", func(t *testing.T) {
		_, err := service.CreateContainer(ctx, NewCreateContainerParams().SetName("web").SetImage("nginx"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only")
	})

	t.Run("RemoveContainer", func(t *testing.T) {
		err := service.RemoveContainer(ctx, NewRemoveContainerParams().SetName("web"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only")
	})

	t.Run("StartContainer", func(t *testing.T) {
		err := service.StartContainer(ctx, NewStartContainerParams().SetName("web"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only")
	})

	t.Run("StopContainer", func(t *testing.T) {
		err := service.StopContainer(ctx, NewStopContainerParams().SetName("web"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only")
	})

	t.Run("RestartContainer", func(t *testing.T) {
		err := service.RestartContainer(ctx, NewRestartContainerParams().SetName("web"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only")
	})
}

func TestService_ReadOnly_AllowsReadOperations(t *testing.T) {
	policy := services.NewPolicy(true)
	mockClient := containermock.NewMockProviderClient(t)
	service := NewService(mockClient, policy, zap.NewNop())

	ctx := context.Background()

	mockClient.On("ListContainers", mock.Anything, mock.Anything).Return([]providers.Container{}, nil)

	_, err := service.ListContainers(ctx, NewListContainersParams())
	require.NoError(t, err)

	mockClient.On("InspectContainer", mock.Anything, mock.Anything).Return(providers.ContainerInspect{}, nil)

	_, err = service.InspectContainer(ctx, NewInspectContainerParams().SetName("web"))
	require.NoError(t, err)
}
