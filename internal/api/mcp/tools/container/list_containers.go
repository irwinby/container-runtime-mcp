package container

import (
	"context"
	"fmt"

	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListContainersInput struct {
	All    bool `json:"all,omitempty" jsonschema:"list all containers, including stopped ones"`
	Limit  int  `json:"limit,omitempty" jsonschema:"limit the number of containers returned"`
	Size   bool `json:"size,omitempty" jsonschema:"include container sizes"`
	Latest bool `json:"latest,omitempty" jsonschema:"show only the latest created container"`
}

type ListContainersItem struct {
	ID      string   `json:"id"`
	Names   []string `json:"names"`
	Image   string   `json:"image"`
	State   string   `json:"state"`
	Status  string   `json:"status"`
	Created int64    `json:"created"`
}

type ListContainersOutput struct {
	Containers []ListContainersItem `json:"containers"`
}

func (i *ListContainersInput) Validate() error {
	if i == nil {
		return fmt.Errorf("list containers input is nil")
	}

	if i.Limit < 0 {
		return fmt.Errorf("limit must be non-negative")
	}

	return nil
}

func (h *ToolsHandler) ListContainers(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input ListContainersInput,
) (*mcp.CallToolResult, ListContainersOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, ListContainersOutput{}, fmt.Errorf("validate list containers input: %w", err)
	}

	params := container.NewListContainersParams().
		SetAll(input.All).
		SetLimit(input.Limit).
		SetSize(input.Size).
		SetLatest(input.Latest)

	containers, err := h.containerService.ListContainers(ctx, params)
	if err != nil {
		return nil, ListContainersOutput{}, fmt.Errorf("list containers: %w", err)
	}

	output := ListContainersOutput{
		Containers: make([]ListContainersItem, 0, len(containers)),
	}

	for _, container := range containers {
		output.Containers = append(output.Containers, ListContainersItem{
			ID:      container.ID,
			Names:   container.Names,
			Image:   container.Image,
			State:   container.State,
			Status:  container.Status,
			Created: container.Created,
		})
	}

	return nil, output, nil
}
