package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Options = mcp.ServerOptions

type Option func(opts *Options)

type Handler interface {
	Register(server *mcp.Server)
}

type Handlers []Handler

// Server runs an MCP server using a configured transport.
type Server interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

func NewHandlers(handlers ...Handler) Handlers {
	return handlers
}
