package container

import (
	"context"
	"errors"
	"testing"

	containermock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/container/mock"
	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerListContainers(t *testing.T) {
	type given struct {
		input  ListContainersInput
		result []container.Container
		err    error
	}

	type want struct {
		called     bool
		containers []ListContainersItem
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				input: ListContainersInput{All: true},
				result: []container.Container{
					{ID: "abc123", Names: []string{"web"}, Image: "nginx:latest", Status: "Up 1 hour"},
				},
			},
			want: want{
				called: true,
				containers: []ListContainersItem{
					{ID: "abc123", Names: []string{"web"}, Image: "nginx:latest", Status: "Up 1 hour"},
				},
			},
		},
		"empty result": {
			given: given{
				input:  ListContainersInput{},
				result: []container.Container{},
			},
			want: want{
				called:     true,
				containers: []ListContainersItem{},
			},
		},
		"service error": {
			given: given{
				input: ListContainersInput{},
				err:   errors.New("docker error"),
			},
			want: want{
				called: true,
			},
		},
		"negative limit": {
			given: given{
				input: ListContainersInput{Limit: -1},
			},
			want: want{
				called: false,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := containermock.NewMockContainerService(t)
			if tt.want.called {
				mockSvc.On("ListContainers", mock.Anything, container.ListContainersParams{
					All:    tt.given.input.All,
					Limit:  tt.given.input.Limit,
					Size:   tt.given.input.Size,
					Latest: tt.given.input.Latest,
				}).Return(tt.given.result, tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, output, err := handler.ListContainers(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			if !tt.want.called {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.containers, output.Containers)
		})
	}
}
