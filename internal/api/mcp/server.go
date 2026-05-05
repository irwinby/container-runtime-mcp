package mcp

import (
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/config"
)

// NewServer creates an MCP server using the configured transport.
func NewServer(cfg config.MCPServer, handlers Handlers, opts ...Option) (Server, error) {
	switch cfg.TransportConfig.Type {
	case config.TransportStdio:
		return NewSTDIOServer(cfg, handlers, opts...), nil
	case config.TransportHTTP:
		return NewHTTPServer(cfg, handlers, opts...)
	default:
		return nil, fmt.Errorf("unsupported transport: %s", cfg.TransportConfig.Type)
	}
}
