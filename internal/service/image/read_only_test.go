package image

import (
	"context"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	services "github.com/irwinby/container-runtime-mcp/internal/service"
	imagemock "github.com/irwinby/container-runtime-mcp/internal/service/image/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestService_ReadOnly(t *testing.T) {
	policy := services.NewPolicy(true)
	mockClient := imagemock.NewMockProviderClient(t)
	service := NewService(mockClient, policy, zap.NewNop())

	ctx := context.Background()

	t.Run("PullImage", func(t *testing.T) {
		err := service.PullImage(ctx, NewPullImageParams().SetRef("nginx:latest"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only")
	})

	t.Run("PushImage", func(t *testing.T) {
		err := service.PushImage(ctx, NewPushImageParams().SetRef("nginx:latest"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only")
	})

	t.Run("RemoveImage", func(t *testing.T) {
		err := service.RemoveImage(ctx, NewRemoveImageParams().SetRef("nginx:latest"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only")
	})

	t.Run("TagImage", func(t *testing.T) {
		err := service.TagImage(ctx, NewTagImageParams().SetSource("nginx:latest").SetTarget("my-nginx:latest"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read-only")
	})
}

func TestService_ReadOnly_AllowsReadOperations(t *testing.T) {
	policy := services.NewPolicy(true)
	mockClient := imagemock.NewMockProviderClient(t)

	service := NewService(mockClient, policy, zap.NewNop())

	ctx := context.Background()

	mockClient.On("ListImages", mock.Anything, mock.Anything).Return([]providers.Image{}, nil)

	_, err := service.ListImages(ctx, NewListImagesParams())
	require.NoError(t, err)

	mockClient.On("InspectImage", mock.Anything, mock.Anything).Return(providers.ImageInspect{}, nil)

	_, err = service.InspectImage(ctx, NewInspectImageParams().SetRef("nginx:latest"))
	require.NoError(t, err)
}
