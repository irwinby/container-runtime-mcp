package telemetry

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/irwinby/container-runtime-mcp/internal/config"
)

type Handler interface {
	Register(mux *http.ServeMux)
}

type Handlers []Handler

// Server serves telemetry endpoints such as liveness, readiness, and pprof.
type Server struct {
	server *http.Server
}

func NewHandler(handlers ...Handler) Handlers {
	return handlers
}

// NewServer creates a telemetry server with the given configuration.
// It returns an error if the configuration is invalid.
func NewServer(cfg config.TelemetryConfig, handlers Handlers) (*Server, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("telemetry is not enabled")
	}

	mux := http.NewServeMux()

	for _, handler := range handlers {
		handler.Register(mux)
	}

	if cfg.PPROFEnabled {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	return &Server{
		server: &http.Server{
			Addr:              cfg.Addr,
			Handler:           mux,
			ReadHeaderTimeout: cfg.ReadTimeout,
			IdleTimeout:       cfg.IDLETimeout,
		},
	}, nil
}

// Run starts the telemetry server and blocks until the context is canceled
// or the server encounters a fatal error.
func (s *Server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.server.Addr, err)
	}

	defer func() {
		// Listener close errors after serve are not critical.
		_ = listener.Close()
	}()

	errs := make(chan error, 1)

	go func() {
		errs <- s.server.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		err := <-errs
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve telemetry: %w", err)
		}

		return nil
	case err := <-errs:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve telemetry: %w", err)
		}

		return nil
	}
}

// Shutdown gracefully shuts down the telemetry server.
func (s *Server) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown telemetry server: %w", err)
	}

	return nil
}
