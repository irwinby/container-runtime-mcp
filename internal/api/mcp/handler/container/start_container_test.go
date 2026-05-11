package container

import (
	"context"
	"errors"
	"testing"

	containermock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/handler/container/mock"
	"github.com/irwinby/container-runtime-mcp/internal/service/container"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerStartContainer(t *testing.T) {
	type given struct {
		input StartContainerInput
		err   error
	}

	type want struct {
		called bool
		name   string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{input: StartContainerInput{Name: "web"}},
			want:  want{called: true, name: "web"},
		},
		"whitespace name": {
			given: given{input: StartContainerInput{Name: "  web  "}},
			want:  want{called: true, name: "web"},
		},
		"empty name": {
			given: given{input: StartContainerInput{Name: ""}, err: errors.New("validation error")},
			want:  want{},
		},
		"service error": {
			given: given{input: StartContainerInput{Name: "web"}, err: errors.New("docker error")},
			want:  want{called: true, name: "web"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := containermock.NewMockContainerService(t)

			if test.want.called {
				mockService.On("StartContainer", mock.Anything, container.StartContainerParams{
					Name: test.want.name,
				}).Return(test.given.err)
			}

			handler := NewToolsHandler(mockService)

			_, _, err := handler.StartContainer(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
