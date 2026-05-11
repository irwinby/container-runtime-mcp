package system

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SystemInfoInput struct{}

type SystemInfo struct {
	ID                string      `json:"id"`
	Containers        int         `json:"containers"`
	ContainersRunning int         `json:"containers_running"`
	ContainersPaused  int         `json:"containers_paused"`
	ContainersStopped int         `json:"containers_stopped"`
	Images            int         `json:"images"`
	Driver            string      `json:"driver"`
	DriverStatus      [][2]string `json:"driver_status"`
	KernelVersion     string      `json:"kernel_version"`
	OperatingSystem   string      `json:"operating_system"`
	OSType            string      `json:"os_type"`
	Architecture      string      `json:"architecture"`
	NCPU              int         `json:"ncpu"`
	MemTotal          int64       `json:"mem_total"`
	ServerVersion     string      `json:"server_version"`
}

type SystemInfoOutput struct {
	Info SystemInfo `json:"info"`
}

func (h *ToolsHandler) SystemInfo(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	_ SystemInfoInput,
) (*mcp.CallToolResult, SystemInfoOutput, error) {
	info, err := h.systemService.SystemInfo(ctx)
	if err != nil {
		return nil, SystemInfoOutput{}, fmt.Errorf("get system info: %w", err)
	}

	return nil, SystemInfoOutput{
		Info: SystemInfo{
			ID:                info.ID,
			Containers:        info.Containers,
			ContainersRunning: info.ContainersRunning,
			ContainersPaused:  info.ContainersPaused,
			ContainersStopped: info.ContainersStopped,
			Images:            info.Images,
			Driver:            info.Driver,
			DriverStatus:      info.DriverStatus,
			KernelVersion:     info.KernelVersion,
			OperatingSystem:   info.OperatingSystem,
			OSType:            info.OSType,
			Architecture:      info.Architecture,
			NCPU:              info.NCPU,
			MemTotal:          info.MemTotal,
			ServerVersion:     info.ServerVersion,
		},
	}, nil
}
