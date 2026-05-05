package mcp

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type STDIOServer struct {
	server *mcp.Server
}

func NewSTDIOServer(cfg config.MCPServer, handlers Handlers, opts ...Option) *STDIOServer {
	serverOpts := &Options{}

	for _, opt := range opts {
		opt(serverOpts)
	}

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    cfg.Name,
			Title:   cfg.Title,
			Version: cfg.Version,
		},
		serverOpts,
	)

	for _, handler := range handlers {
		handler.Register(server)
	}

	return &STDIOServer{
		server: server,
	}
}

func (s *STDIOServer) Run(ctx context.Context) error {
	err := s.server.Run(ctx, &mcp.StdioTransport{})
	if err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}

func (s *STDIOServer) Shutdown(ctx context.Context) error {
	return nil
}
