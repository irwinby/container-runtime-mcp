package mcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mcpmock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/mock"
	"github.com/irwinby/container-runtime-mcp/internal/config"
	testnet "github.com/irwinby/container-runtime-mcp/internal/testing/net"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewHandlers(t *testing.T) {
	h1 := mcpmock.NewMockHandler(t)
	h2 := mcpmock.NewMockHandler(t)

	handlers := NewHandlers(h1, h2)

	assert.Len(t, handlers, 2)
	assert.Equal(t, h1, handlers[0])
	assert.Equal(t, h2, handlers[1])
}

func TestNewServer(t *testing.T) {
	type given struct {
		cfg      config.MCPServer
		handlers Handlers
		opts     []Option
	}

	type want struct {
		registered bool
		server     any
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"server with no handlers": {
			given: given{
				cfg: config.MCPServer{
					Name:    "Test",
					Version: "1.0.0",
					TransportConfig: config.TransportConfig{
						Type: config.TransportStdio,
					},
				},
			},
			want: want{
				registered: false,
				server:     &STDIOServer{},
			},
		},
		"server with handlers": {
			given: given{
				cfg: config.MCPServer{
					Name:    "Test",
					Version: "1.0.0",
					TransportConfig: config.TransportConfig{
						Type: config.TransportStdio,
					},
				},
			},
			want: want{
				registered: true,
				server:     &STDIOServer{},
			},
		},
		"server with options": {
			given: given{
				cfg: config.MCPServer{
					Name:    "Test",
					Version: "1.0.0",
					TransportConfig: config.TransportConfig{
						Type: config.TransportStdio,
					},
				},
				opts: []Option{
					func(opts *Options) {
						opts.Capabilities = &mcp.ServerCapabilities{}
					},
				},
			},
			want: want{
				registered: true,
				server:     &STDIOServer{},
			},
		},
		"http server": {
			given: given{
				cfg: config.MCPServer{
					Name:    "Test",
					Version: "1.0.0",
					TransportConfig: config.TransportConfig{
						Type: config.TransportHTTP,
						HTTP: &config.HTTPTransportConfig{
							Addr:           "127.0.0.1:8080",
							Path:           "/mcp",
							SessionTimeout: 0,
						},
					},
				},
			},
			want: want{
				registered: false,
				server:     &HTTPServer{},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			given := test.given

			if test.want.registered {
				mockHandler := mcpmock.NewMockHandler(t)

				mockHandler.On("Register", mock.Anything).Once()

				given.handlers = NewHandlers(mockHandler)
			}

			server, err := NewServer(given.cfg, given.handlers, given.opts...)

			require.NoError(t, err)
			assert.NotNil(t, server)
			assert.IsType(t, test.want.server, server)
		})
	}
}

func TestNewServer_UnsupportedTransport(t *testing.T) {
	server, err := NewServer(
		config.MCPServer{
			Name:            "Test",
			Version:         "1.0.0",
			TransportConfig: config.TransportConfig{Type: "ws"},
		},
		nil,
	)

	require.Error(t, err)
	assert.Nil(t, server)
	assert.Contains(t, err.Error(), "unsupported transport")
}

func TestServer_HTTPRouting_ExactPathOnly(t *testing.T) {
	mockH := mcpmock.NewMockHandler(t)
	mockH.On("Register", mock.Anything).Once()

	server, err := NewServer(
		config.MCPServer{
			Name:    "Test",
			Version: "1.0.0",
			TransportConfig: config.TransportConfig{
				Type: config.TransportHTTP,
				HTTP: &config.HTTPTransportConfig{
					Addr:           "127.0.0.1:8080",
					Path:           "/mcp",
					SessionTimeout: 0,
				},
			},
		},
		NewHandlers(mockH),
	)
	require.NoError(t, err)

	httpServer, ok := server.(*HTTPServer)
	require.True(t, ok)

	mux, ok := httpServer.server.Handler.(*http.ServeMux)
	require.True(t, ok)

	// Exact path is registered and reaches the handler.
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	mux.ServeHTTP(recorder, request)
	assert.NotEqual(t, http.StatusNotFound, recorder.Code, "exact path should be routed")

	// Subpath is not registered and returns 404.
	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/mcp/extra", nil)
	mux.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNotFound, recorder.Code, "subpath should return 404")

	// Trailing slash is not registered and returns 404.
	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/mcp/", nil)
	mux.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNotFound, recorder.Code, "trailing slash should return 404")
}

func TestServer_HTTPRouting_AuthRequired(t *testing.T) {
	mockHandler := mcpmock.NewMockHandler(t)

	mockHandler.On("Register", mock.Anything).Once()

	server, err := NewServer(
		config.MCPServer{
			Name:    "Test",
			Version: "1.0.0",
			TransportConfig: config.TransportConfig{
				Type: config.TransportHTTP,
				HTTP: &config.HTTPTransportConfig{
					Addr:           "127.0.0.1:8080",
					Path:           "/mcp",
					SessionTimeout: 0,
					AuthToken:      "test-token",
				},
			},
		},
		NewHandlers(mockHandler),
	)

	require.NoError(t, err)

	httpServer, ok := server.(*HTTPServer)
	require.True(t, ok)

	mux, ok := httpServer.server.Handler.(*http.ServeMux)
	require.True(t, ok)

	// Missing Authorization header returns 401.
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	mux.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, "missing auth should return 401")

	// Wrong token returns 401.
	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/mcp", nil)
	request.Header.Set("Authorization", "Bearer wrong-token")
	mux.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, "wrong token should return 401")

	// Correct token reaches the handler.
	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/mcp", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	mux.ServeHTTP(recorder, request)
	assert.NotEqual(t, http.StatusUnauthorized, recorder.Code, "correct token should not return 401")
	assert.NotEqual(t, http.StatusNotFound, recorder.Code, "exact path should be routed")
}

func TestHTTPServer_Shutdown(t *testing.T) {
	mockH := mcpmock.NewMockHandler(t)
	mockH.On("Register", mock.Anything).Once()

	addr := testnet.FreeTCPAddr(t)

	server, err := NewServer(
		config.MCPServer{
			Name:    "Test",
			Version: "1.0.0",
			TransportConfig: config.TransportConfig{
				Type: config.TransportHTTP,
				HTTP: &config.HTTPTransportConfig{
					Addr:           addr,
					Path:           "/mcp",
					SessionTimeout: 0,
				},
			},
		},
		NewHandlers(mockH),
	)
	require.NoError(t, err)

	httpServer, ok := server.(*HTTPServer)
	require.True(t, ok)

	errs := make(chan error, 1)

	go func() {
		errs <- httpServer.Run(context.Background())
	}()

	// Wait for the server to start listening.
	testnet.RequireTCPListening(t, addr)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = httpServer.Shutdown(shutdownCtx)
	require.NoError(t, err)

	runErr := <-errs
	require.NoError(t, runErr)
}

func TestHTTPServer_Run_ContextCancelled(t *testing.T) {
	mockHandler := mcpmock.NewMockHandler(t)

	mockHandler.On("Register", mock.Anything).Once()

	addr := testnet.FreeTCPAddr(t)

	server, err := NewServer(
		config.MCPServer{
			Name:    "Test",
			Version: "1.0.0",
			TransportConfig: config.TransportConfig{
				Type: config.TransportHTTP,
				HTTP: &config.HTTPTransportConfig{
					Addr:           addr,
					Path:           "/mcp",
					SessionTimeout: 0,
				},
			},
		},
		NewHandlers(mockHandler),
	)
	require.NoError(t, err)

	httpServer, ok := server.(*HTTPServer)
	require.True(t, ok)

	ctx, cancel := context.WithCancel(context.Background())

	errs := make(chan error, 1)

	go func() {
		errs <- httpServer.Run(ctx)
	}()

	// Wait for the server to start listening.
	testnet.RequireTCPListening(t, addr)

	cancel()

	// Run waits for the server to finish after context cancellation.
	// Trigger shutdown so Serve returns.
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	err = httpServer.Shutdown(shutdownCtx)
	require.NoError(t, err)

	runErr := <-errs
	require.NoError(t, runErr)
}

func TestHTTPServer_Run_ListenError(t *testing.T) {
	server := &HTTPServer{
		server: &http.Server{Addr: "invalid:address:too:many:colons"},
	}

	err := server.Run(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "listen")
}

func TestHTTPServer_Shutdown_AlreadyClosed(t *testing.T) {
	mockH := mcpmock.NewMockHandler(t)
	mockH.On("Register", mock.Anything).Once()

	addr := testnet.FreeTCPAddr(t)

	server, err := NewServer(
		config.MCPServer{
			Name:    "Test",
			Version: "1.0.0",
			TransportConfig: config.TransportConfig{
				Type: config.TransportHTTP,
				HTTP: &config.HTTPTransportConfig{
					Addr:           addr,
					Path:           "/mcp",
					SessionTimeout: 0,
				},
			},
		},
		NewHandlers(mockH),
	)
	require.NoError(t, err)

	httpServer, ok := server.(*HTTPServer)
	require.True(t, ok)

	err = httpServer.Shutdown(context.Background())
	// Shutdown on a non-started server may return nil or ErrServerClosed.
	// Either is acceptable.
	if err != nil {
		assert.ErrorIs(t, err, http.ErrServerClosed)
	}
}

func TestNewServer_HTTPWithSessionTimeout(t *testing.T) {
	mockH := mcpmock.NewMockHandler(t)
	mockH.On("Register", mock.Anything).Once()

	server, err := NewServer(
		config.MCPServer{
			Name:    "Test",
			Version: "1.0.0",
			TransportConfig: config.TransportConfig{
				Type: config.TransportHTTP,
				HTTP: &config.HTTPTransportConfig{
					Addr:           "127.0.0.1:8080",
					Path:           "/mcp",
					SessionTimeout: 30 * time.Minute,
				},
			},
		},
		NewHandlers(mockH),
	)
	require.NoError(t, err)
	require.NotNil(t, server)
}

func TestNewServer_HTTPNilConfig(t *testing.T) {
	server, err := NewServer(
		config.MCPServer{
			Name:    "Test",
			Version: "1.0.0",
			TransportConfig: config.TransportConfig{
				Type: config.TransportHTTP,
				HTTP: nil,
			},
		},
		nil,
	)

	require.Error(t, err)
	assert.Nil(t, server)
	assert.Contains(t, err.Error(), "http transport configuration is nil")
}
