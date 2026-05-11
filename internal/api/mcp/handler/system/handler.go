package system

import (
	"context"

	"github.com/irwinby/container-runtime-mcp/internal/service/system"
	"github.com/irwinby/container-runtime-mcp/pkg/ptr"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type systemService interface {
	SystemInfo(ctx context.Context) (system.SystemInfo, error)
	SystemVersion(ctx context.Context) (system.SystemVersion, error)
	Ping(ctx context.Context) (system.PingResult, error)
}

type ToolsHandler struct {
	systemService systemService
}

func NewToolsHandler(systemService systemService) *ToolsHandler {
	return &ToolsHandler{
		systemService: systemService,
	}
}

func (h *ToolsHandler) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "info",
		Title:       "Docker Info",
		Description: "Get Docker system information.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.SystemInfo)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "version",
		Title:       "Docker Version",
		Description: "Get Docker version information.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.Version)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "ping",
		Title:       "Ping Docker",
		Description: "Ping the Docker daemon to check connectivity.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.Ping)
}
