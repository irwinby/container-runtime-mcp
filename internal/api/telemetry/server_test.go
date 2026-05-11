package telemetry

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/irwinby/container-runtime-mcp/internal/api/telemetry/handler/probe"
	"github.com/irwinby/container-runtime-mcp/internal/config"
	systemservice "github.com/irwinby/container-runtime-mcp/internal/service/system"
	testnet "github.com/irwinby/container-runtime-mcp/internal/testing/net"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHealthChecker struct {
	err error
}

func (m *mockHealthChecker) Ping(ctx context.Context) (systemservice.PingResult, error) {
	return systemservice.PingResult{}, m.err
}

func TestNewServer_NotEnabled(t *testing.T) {
	cfg := config.TelemetryConfig{Enabled: false}
	checker := &mockHealthChecker{}

	server, err := NewServer(cfg, NewHandler(probe.NewHandler(checker)))
	require.Error(t, err)

	assert.Nil(t, server)
	assert.Contains(t, err.Error(), "telemetry is not enabled")
}

func TestServer_Livez(t *testing.T) {
	cfg := config.TelemetryConfig{
		Enabled: true,
		Addr:    "127.0.0.1:0",
	}

	checker := &mockHealthChecker{}

	server, err := NewServer(cfg, NewHandler(probe.NewHandler(checker)))
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/livez", nil)

	server.server.Handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestServer_Ready_Healthy(t *testing.T) {
	cfg := config.TelemetryConfig{
		Enabled: true,
		Addr:    "127.0.0.1:0",
	}

	checker := &mockHealthChecker{err: nil}

	server, err := NewServer(cfg, NewHandler(probe.NewHandler(checker)))
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)

	server.server.Handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestServer_Ready_Unhealthy(t *testing.T) {
	cfg := config.TelemetryConfig{
		Enabled: true,
		Addr:    "127.0.0.1:0",
	}

	checker := &mockHealthChecker{err: errors.New("ping failed")}

	server, err := NewServer(cfg, NewHandler(probe.NewHandler(checker)))
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)

	server.server.Handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
}

func TestServer_Pprof_Disabled(t *testing.T) {
	cfg := config.TelemetryConfig{
		Enabled:      true,
		Addr:         "127.0.0.1:0",
		PPROFEnabled: false,
	}

	checker := &mockHealthChecker{}

	server, err := NewServer(cfg, NewHandler(probe.NewHandler(checker)))
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)

	server.server.Handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestServer_Pprof_Enabled(t *testing.T) {
	cfg := config.TelemetryConfig{
		Enabled:      true,
		Addr:         "127.0.0.1:0",
		PPROFEnabled: true,
	}

	checker := &mockHealthChecker{}

	server, err := NewServer(cfg, NewHandler(probe.NewHandler(checker)))
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)

	server.server.Handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestServer_Run_ListenError(t *testing.T) {
	server := &Server{
		server: &http.Server{Addr: "invalid:address:too:many:colons"},
	}

	err := server.Run(context.Background())
	require.Error(t, err)

	assert.Contains(t, err.Error(), "listen")
}

func TestServer_Run_AndShutdown(t *testing.T) {
	addr := testnet.FreeTCPAddr(t)

	cfg := config.TelemetryConfig{
		Enabled: true,
		Addr:    addr,
	}

	checker := &mockHealthChecker{}

	server, err := NewServer(cfg, NewHandler(probe.NewHandler(checker)))
	require.NoError(t, err)

	errs := make(chan error, 1)

	go func() {
		errs <- server.Run(context.Background())
	}()

	testnet.RequireTCPListening(t, addr)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	require.NoError(t, err)

	runErr := <-errs
	require.NoError(t, runErr)
}

func TestServer_Run_ContextCancelled(t *testing.T) {
	addr := testnet.FreeTCPAddr(t)

	cfg := config.TelemetryConfig{
		Enabled: true,
		Addr:    addr,
	}

	checker := &mockHealthChecker{}

	server, err := NewServer(cfg, NewHandler(probe.NewHandler(checker)))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	errs := make(chan error, 1)

	go func() {
		errs <- server.Run(ctx)
	}()

	testnet.RequireTCPListening(t, addr)

	cancel()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	err = server.Shutdown(shutdownCtx)
	require.NoError(t, err)

	runErr := <-errs
	require.NoError(t, runErr)
}

func TestServer_Shutdown_AlreadyClosed(t *testing.T) {
	cfg := config.TelemetryConfig{
		Enabled: true,
		Addr:    "127.0.0.1:0",
	}

	checker := &mockHealthChecker{}

	server, err := NewServer(cfg, NewHandler(probe.NewHandler(checker)))
	require.NoError(t, err)

	err = server.Shutdown(context.Background())
	// Shutdown on a non-started server may return nil or ErrServerClosed.
	// Either is acceptable.
	if err != nil {
		assert.ErrorIs(t, err, http.ErrServerClosed)
	}
}

func TestServer_UnknownPath(t *testing.T) {
	cfg := config.TelemetryConfig{
		Enabled: true,
		Addr:    "127.0.0.1:0",
	}

	checker := &mockHealthChecker{}

	server, err := NewServer(cfg, NewHandler(probe.NewHandler(checker)))
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/unknown", nil)

	server.server.Handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}
