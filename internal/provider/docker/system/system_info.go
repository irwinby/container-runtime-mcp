package system

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) SystemInfo(ctx context.Context) (providers.SystemInfo, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()
	result, err := p.client.Info(ctx, client.InfoOptions{})
	if err != nil {
		return providers.SystemInfo{}, fmt.Errorf("get docker system info: %w", err)
	}

	info := result.Info

	return providers.NewSystemInfo().
		SetID(info.ID).
		SetContainers(info.Containers).
		SetContainersRunning(info.ContainersRunning).
		SetContainersPaused(info.ContainersPaused).
		SetContainersStopped(info.ContainersStopped).
		SetImages(info.Images).
		SetDriver(info.Driver).
		SetDriverStatus(info.DriverStatus).
		SetKernelVersion(info.KernelVersion).
		SetOperatingSystem(info.OperatingSystem).
		SetOSType(info.OSType).
		SetArchitecture(info.Architecture).
		SetNCPU(info.NCPU).
		SetMemTotal(info.MemTotal).
		SetServerVersion(info.ServerVersion), nil
}
