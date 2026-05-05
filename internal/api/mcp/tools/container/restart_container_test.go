package container

import (
	"context"
	"errors"
	"testing"

	containermock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/container/mock"
	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerRestartContainer(t *testing.T) {
	type given struct {
		input RestartContainerInput
		err   error
	}

	type want struct {
		called bool
		name   string
		signal string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{input: RestartContainerInput{Name: "web"}},
			want:  want{called: true, name: "web"},
		},
		"whitespace name": {
			given: given{input: RestartContainerInput{Name: "  web  "}},
			want:  want{called: true, name: "web"},
		},
		"with signal": {
			given: given{input: RestartContainerInput{Name: "web", Signal: "SIGKILL"}},
			want:  want{called: true, name: "web", signal: "SIGKILL"},
		},
		"empty name": {
			given: given{input: RestartContainerInput{Name: ""}, err: errors.New("validation error")},
			want:  want{},
		},
		"service error": {
			given: given{input: RestartContainerInput{Name: "web"}, err: errors.New("docker error")},
			want:  want{called: true, name: "web"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := containermock.NewMockContainerService(t)
			if tt.want.called {
				mockSvc.On("RestartContainer", mock.Anything, container.RestartContainerParams{
					Name:   tt.want.name,
					Signal: tt.want.signal,
				}).Return(tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, _, err := handler.RestartContainer(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
