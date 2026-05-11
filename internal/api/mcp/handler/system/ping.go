package system

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type PingInput struct{}

type PingResult struct {
	APIVersion     string `json:"api_version"`
	OSType         string `json:"os_type"`
	Experimental   bool   `json:"experimental"`
	BuilderVersion string `json:"builder_version"`
}

type PingOutput struct {
	Ping PingResult `json:"ping"`
}

func (h *ToolsHandler) Ping(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	_ PingInput,
) (*mcp.CallToolResult, PingOutput, error) {
	ping, err := h.systemService.Ping(ctx)
	if err != nil {
		return nil, PingOutput{}, fmt.Errorf("ping: %w", err)
	}

	return nil, PingOutput{
		Ping: PingResult{
			APIVersion:     ping.APIVersion,
			OSType:         ping.OSType,
			Experimental:   ping.Experimental,
			BuilderVersion: ping.BuilderVersion,
		},
	}, nil
}
