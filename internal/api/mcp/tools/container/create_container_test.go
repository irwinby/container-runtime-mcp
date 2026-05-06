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

func TestHandlerCreateContainer(t *testing.T) {
	type given struct {
		input CreateContainerInput
		id    string
		err   error
	}

	type want struct {
		called bool
		name   string
		image  string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				input: CreateContainerInput{Name: "web", Image: "nginx:latest"},
				id:    "abc123",
			},
			want: want{
				called: true,
				name:   "web",
				image:  "nginx:latest",
			},
		},
		"whitespace input": {
			given: given{
				input: CreateContainerInput{Name: "  web  ", Image: "  nginx:latest  "},
				id:    "abc123",
			},
			want: want{
				called: true,
				name:   "web",
				image:  "nginx:latest",
			},
		},
		"empty name": {
			given: given{
				input: CreateContainerInput{Name: "", Image: "nginx:latest"},
				err:   errors.New("validation error"),
			},
			want: want{},
		},
		"empty image": {
			given: given{
				input: CreateContainerInput{Name: "web", Image: ""},
				err:   errors.New("validation error"),
			},
			want: want{},
		},
		"service error": {
			given: given{
				input: CreateContainerInput{Name: "web", Image: "nginx:latest"},
				err:   errors.New("docker error"),
			},
			want: want{
				called: true,
				name:   "web",
				image:  "nginx:latest",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := containermock.NewMockContainerService(t)

			if test.want.called {
				mockService.On("CreateContainer", mock.Anything, container.CreateContainerParams{
					Name:  test.want.name,
					Image: test.want.image,
				}).Return(test.given.id, test.given.err)
			}

			handler := NewToolsHandler(mockService)

			_, output, err := handler.CreateContainer(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.given.id, output.ID)
		})
	}
}
