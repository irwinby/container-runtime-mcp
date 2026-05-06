package mcp

import (
	"context"
	"testing"
	"time"

	mcpmock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/mock"
	"github.com/irwinby/container-runtime-mcp/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewSTDIOServer(t *testing.T) {
	mockH := mcpmock.NewMockHandler(t)
	mockH.On("Register", mock.Anything).Once()

	server := NewSTDIOServer(
		config.MCPServer{Name: "Test", Version: "1.0.0"},
		NewHandlers(mockH),
		func(opts *Options) {
			opts.Capabilities = &mcp.ServerCapabilities{}
		},
	)

	require.NotNil(t, server)
	assert.NotNil(t, server.server)
}

func TestSTDIOServer_Run(t *testing.T) {
	mockH := mcpmock.NewMockHandler(t)
	mockH.On("Register", mock.Anything).Once()

	server := NewSTDIOServer(
		config.MCPServer{Name: "Test", Version: "1.0.0"},
		NewHandlers(mockH),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// Cancel after a short delay to let Run start.
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	err := server.Run(ctx)
	// Run may return context.Canceled when the context is canceled;
	// that is acceptable.
	if err != nil {
		require.ErrorIs(t, err, context.Canceled)
	}
}

func TestSTDIOServer_Shutdown(t *testing.T) {
	mockH := mcpmock.NewMockHandler(t)
	mockH.On("Register", mock.Anything).Once()

	server := NewSTDIOServer(
		config.MCPServer{Name: "Test", Version: "1.0.0"},
		NewHandlers(mockH),
	)

	err := server.Shutdown(context.Background())
	require.NoError(t, err)
}
