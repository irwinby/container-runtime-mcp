package container

import (
	"context"
	"errors"
	"testing"

	containermock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/handler/container/mock"
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
		"nil input": {
			given: given{
				input: ContainerLogsInput{},
			},
			want: want{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := containermock.NewMockContainerService(t)

			if test.want.called {
				mockService.On("ContainerLogs", mock.Anything, container.ContainerLogsParams{
					Name:   test.want.name,
					Stdout: true,
					Stderr: true,
					Since:  test.given.input.Since,
					Tail:   test.given.input.Tail,
				}).Return(test.given.result, test.given.err)
			}

			handler := NewToolsHandler(mockService)

			_, output, err := handler.ContainerLogs(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			if !test.want.called {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want.logs, output)
		})
	}
}
