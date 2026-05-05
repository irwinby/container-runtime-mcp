package system

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) SystemVersion(ctx context.Context) (providers.SystemVersion, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()
	result, err := p.client.ServerVersion(ctx, client.ServerVersionOptions{})
	if err != nil {
		return providers.SystemVersion{}, fmt.Errorf("get docker system version: %w", err)
	}

	return providers.NewSystemVersion().
		SetVersion(result.Version).
		SetAPIVersion(result.APIVersion).
		SetMinAPIVersion(result.MinAPIVersion).
		SetOs(result.Os).
		SetArch(result.Arch).
		SetPlatformName(result.Platform.Name), nil
}
