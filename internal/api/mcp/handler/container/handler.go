package container

import (
	"context"

	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/irwinby/container-runtime-mcp/pkg/ptr"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type containerService interface {
	CanWrite() bool
	CreateContainer(ctx context.Context, params container.CreateContainerParams) (string, error)
	RemoveContainer(ctx context.Context, params container.RemoveContainerParams) error
	ListContainers(ctx context.Context, params container.ListContainersParams) ([]container.Container, error)
	InspectContainer(ctx context.Context, params container.InspectContainerParams) (container.ContainerInspect, error)
	StartContainer(ctx context.Context, params container.StartContainerParams) error
	StopContainer(ctx context.Context, params container.StopContainerParams) error
	RestartContainer(ctx context.Context, params container.RestartContainerParams) error
	ContainerLogs(ctx context.Context, params container.ContainerLogsParams) (container.ContainerLogsResult, error)
	ExecContainer(ctx context.Context, params container.ExecContainerParams) (container.ExecContainerResult, error)
}

type ToolsHandler struct {
	containerService containerService
}

func NewToolsHandler(containerService containerService) *ToolsHandler {
	return &ToolsHandler{
		containerService: containerService,
	}
}

func (h *ToolsHandler) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_containers",
		Title:       "List Containers",
		Description: "List Docker containers.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.ListContainers)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "inspect_container",
		Title:       "Inspect Container",
		Description: "Inspect a Docker container.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.InspectContainer)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "container_logs",
		Title:       "Container Logs",
		Description: "Get logs from a Docker container.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:  true,
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.ContainerLogs)

	if !h.containerService.CanWrite() {
		return
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_container",
		Title:       "Create Container",
		Description: "Create a new Docker container.",
		Annotations: &mcp.ToolAnnotations{
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.CreateContainer)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "remove_container",
		Title:       "Remove Container",
		Description: "Remove a Docker container.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: ptr.Bool(true),
			OpenWorldHint:   ptr.Bool(false),
		},
	}, h.RemoveContainer)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "start_container",
		Title:       "Start Container",
		Description: "Start a Docker container.",
		Annotations: &mcp.ToolAnnotations{
			OpenWorldHint: ptr.Bool(false),
		},
	}, h.StartContainer)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "stop_container",
		Title:       "Stop Container",
		Description: "Stop a Docker container.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: ptr.Bool(true),
			OpenWorldHint:   ptr.Bool(false),
		},
	}, h.StopContainer)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "restart_container",
		Title:       "Restart Container",
		Description: "Restart a Docker container.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: ptr.Bool(true),
			OpenWorldHint:   ptr.Bool(false),
		},
	}, h.RestartContainer)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "exec_container",
		Title:       "Exec Container",
		Description: "Execute a command in a running Docker container.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: ptr.Bool(true),
			OpenWorldHint:   ptr.Bool(false),
		},
	}, h.ExecContainer)
}
