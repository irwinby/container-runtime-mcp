package mcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/irwinby/container-runtime-mcp/internal/api/mcp/middleware"
	"github.com/irwinby/container-runtime-mcp/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type HTTPServer struct {
	server *http.Server
}

// NewHTTPServer creates an MCP server over HTTP streamable transport.
// It validates the HTTP configuration and returns an error for invalid input.
func NewHTTPServer(cfg config.MCPServer, handlers Handlers, opts ...Option) (*HTTPServer, error) {
	if cfg.TransportConfig.HTTP == nil {
		return nil, fmt.Errorf("http transport configuration is nil")
	}

	serverOpts := &Options{}

	for _, opt := range opts {
		opt(serverOpts)
	}

	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    cfg.Name,
			Title:   cfg.Title,
			Version: cfg.Version,
		},
		serverOpts,
	)

	for _, handler := range handlers {
		handler.Register(mcpServer)
	}

	streamOpts := &mcp.StreamableHTTPOptions{}

	if cfg.TransportConfig.HTTP.SessionTimeout > 0 {
		streamOpts.SessionTimeout = cfg.TransportConfig.HTTP.SessionTimeout
	}

	handler := http.Handler(mcp.NewStreamableHTTPHandler(
		func(_ *http.Request) *mcp.Server {
			return mcpServer
		},
		streamOpts,
	))

	mux := http.NewServeMux()
	mux.Handle(cfg.TransportConfig.HTTP.Path, middleware.Auth(cfg.TransportConfig.HTTP.AuthToken, handler))

	return &HTTPServer{
		server: &http.Server{
			Addr:              cfg.TransportConfig.HTTP.Addr,
			Handler:           mux,
			ReadHeaderTimeout: cfg.TransportConfig.HTTP.ReadTimeout,
			IdleTimeout:       cfg.TransportConfig.HTTP.IDLETimeout,
		},
	}, nil
}

func (s *HTTPServer) Run(ctx context.Context) error {
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
		// Wait for the server to finish (caller is responsible for shutdown).
		err := <-errs
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve http: %w", err)
		}

		return nil
	case err := <-errs:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve http: %w", err)
		}

		return nil
	}
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	return nil
}
