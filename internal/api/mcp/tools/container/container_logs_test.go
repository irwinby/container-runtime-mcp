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

func TestHandlerContainerLogs(t *testing.T) {
	type given struct {
		input  ContainerLogsInput
		result container.ContainerLogsResult
		err    error
	}

	type want struct {
		called bool
		name   string
		logs   ContainerLogsOutput
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				input:  ContainerLogsInput{Name: "web"},
				result: container.ContainerLogsResult{Stdout: "hello", Stderr: "warn"},
			},
			want: want{
				called: true,
				name:   "web",
				logs:   ContainerLogsOutput{Stdout: "hello", Stderr: "warn"},
			},
		},
		"empty name": {
			given: given{
				input: ContainerLogsInput{Name: ""},
				err:   errors.New("validation error"),
			},
			want: want{},
		},
		"service error": {
			given: given{
				input: ContainerLogsInput{Name: "web"},
				err:   errors.New("docker error"),
			},
			want: want{
				called: true,
				name:   "web",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := containermock.NewMockContainerService(t)
			if tt.want.called {
				mockSvc.On("ContainerLogs", mock.Anything, container.ContainerLogsParams{
					Name:   tt.want.name,
					Stdout: true,
					Stderr: true,
					Since:  tt.given.input.Since,
					Tail:   tt.given.input.Tail,
				}).Return(tt.given.result, tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, output, err := handler.ContainerLogs(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			if !tt.want.called {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.logs, output)
		})
	}
}
