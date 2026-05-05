package system

import (
	"context"
	"fmt"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	"github.com/moby/moby/client"
)

func (p *Provider) Ping(ctx context.Context) (providers.PingResult, error) {
	ctx, cancel := p.withTimeout(ctx)
	defer cancel()
	result, err := p.client.Ping(ctx, client.PingOptions{})
	if err != nil {
		return providers.PingResult{}, fmt.Errorf("ping docker: %w", err)
	}

	return providers.NewPingResult().
		SetAPIVersion(result.APIVersion).
		SetOSType(result.OSType).
		SetExperimental(result.Experimental).
		SetBuilderVersion(string(result.BuilderVersion)), nil
}
