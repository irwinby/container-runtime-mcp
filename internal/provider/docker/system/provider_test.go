package system

import (
	"context"
	"testing"
	"time"

	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func nopTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return ctx, func() {}
}

func TestNewProvider(t *testing.T) {
	mockClient := dockermock.NewMockDockerClient(t)
	provider := NewProvider(mockClient, nopTimeout)
	require.NotNil(t, provider)
	assert.Equal(t, mockClient, provider.client)
}

func TestProviderWithTimeout(t *testing.T) {
	mockClient := dockermock.NewMockDockerClient(t)

	called := false
	withTimeout := func(ctx context.Context) (context.Context, context.CancelFunc) {
		called = true
		return ctx, func() {}
	}

	provider := NewProvider(mockClient, withTimeout)
	_, cancel := provider.WithTimeout(context.Background())
	cancel()
	require.True(t, called)
}

func TestProviderWithTimeout_RealTimeout(t *testing.T) {
	mockClient := dockermock.NewMockDockerClient(t)

	withTimeout := func(ctx context.Context) (context.Context, context.CancelFunc) {
		return context.WithTimeout(ctx, time.Millisecond)
	}

	provider := NewProvider(mockClient, withTimeout)
	ctx, cancel := provider.WithTimeout(context.Background())
	defer cancel()
	require.NotNil(t, ctx)
}
