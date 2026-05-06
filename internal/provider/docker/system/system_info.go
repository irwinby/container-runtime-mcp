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

	return providers.NewSystemInfo().
		SetID(result.Info.ID).
		SetContainers(result.Info.Containers).
		SetContainersRunning(result.Info.ContainersRunning).
		SetContainersPaused(result.Info.ContainersPaused).
		SetContainersStopped(result.Info.ContainersStopped).
		SetImages(result.Info.Images).
		SetDriver(result.Info.Driver).
		SetDriverStatus(result.Info.DriverStatus).
		SetKernelVersion(result.Info.KernelVersion).
		SetOperatingSystem(result.Info.OperatingSystem).
		SetOSType(result.Info.OSType).
		SetArchitecture(result.Info.Architecture).
		SetNCPU(result.Info.NCPU).
		SetMemTotal(result.Info.MemTotal).
		SetServerVersion(result.Info.ServerVersion), nil
}
