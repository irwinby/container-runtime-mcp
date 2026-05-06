package docker

import (
	"context"
	"testing"
	"time"

	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProvider(t *testing.T) {
	ctx := context.Background()
	provider, err := NewProvider(ctx, time.Minute)

	if err != nil {
		// Docker may not be available in the test environment.
		require.Contains(t, err.Error(), "create docker client")
		return
	}

	require.NotNil(t, provider)
	require.NoError(t, provider.Close())
}

func TestNewProvider_ZeroTimeout(t *testing.T) {
	mockClient := dockermock.NewMockDockerClient(t)
	provider := newProvider(mockClient, 0)

	require.NotNil(t, provider)
	assert.NotNil(t, provider.ContainerProvider)
	assert.NotNil(t, provider.ImageProvider)
	assert.NotNil(t, provider.VolumeProvider)
	assert.NotNil(t, provider.SystemProvider)
	assert.Equal(t, mockClient, provider.client)
	assert.Equal(t, time.Duration(0), provider.timeout)
}

func TestNewProvider_WithTimeout(t *testing.T) {
	mockClient := dockermock.NewMockDockerClient(t)
	provider := newProvider(mockClient, time.Minute)

	require.NotNil(t, provider)
	assert.NotNil(t, provider.ContainerProvider)
	assert.NotNil(t, provider.ImageProvider)
	assert.NotNil(t, provider.VolumeProvider)
	assert.NotNil(t, provider.SystemProvider)
	assert.Equal(t, mockClient, provider.client)
	assert.Equal(t, time.Minute, provider.timeout)
}

func TestProviderClose(t *testing.T) {
	mockClient := dockermock.NewMockDockerClient(t)
	mockClient.On("Close").Return(nil).Once()

	provider := newProvider(mockClient, 0)
	err := provider.Close()
	require.NoError(t, err)
}

func TestProviderClose_Error(t *testing.T) {
	mockClient := dockermock.NewMockDockerClient(t)
	mockClient.On("Close").Return(assert.AnError).Once()

	provider := newProvider(mockClient, 0)
	err := provider.Close()
	require.Error(t, err)
}
