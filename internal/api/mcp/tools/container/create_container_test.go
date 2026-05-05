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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := containermock.NewMockContainerService(t)
			if tt.want.called {
				mockSvc.On("CreateContainer", mock.Anything, container.CreateContainerParams{
					Name:  tt.want.name,
					Image: tt.want.image,
				}).Return(tt.given.id, tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, output, err := handler.CreateContainer(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.given.id, output.ID)
		})
	}
}
