package system

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type VersionInput struct{}

type Version struct {
	Version       string `json:"version"`
	APIVersion    string `json:"api_version"`
	MinAPIVersion string `json:"min_api_version"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	PlatformName  string `json:"platform_name"`
}

type VersionOutput struct {
	Version Version `json:"version"`
}

func (h *ToolsHandler) Version(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	_ VersionInput,
) (*mcp.CallToolResult, VersionOutput, error) {
	version, err := h.systemService.SystemVersion(ctx)
	if err != nil {
		return nil, VersionOutput{}, fmt.Errorf("get version: %w", err)
	}

	return nil, VersionOutput{
		Version: Version{
			Version:       version.Version,
			APIVersion:    version.APIVersion,
			MinAPIVersion: version.MinAPIVersion,
			OS:            version.Os,
			Arch:          version.Arch,
			PlatformName:  version.PlatformName,
		},
	}, nil
}
