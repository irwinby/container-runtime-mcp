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

func TestHandlerInspectContainer(t *testing.T) {
	type given struct {
		input  InspectContainerInput
		result container.ContainerInspect
		err    error
	}

	type want struct {
		called    bool
		name      string
		container InspectContainerDetails
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				input:  InspectContainerInput{Name: "web"},
				result: container.ContainerInspect{ID: "abc", Name: "web", State: "running"},
			},
			want: want{
				called:    true,
				name:      "web",
				container: InspectContainerDetails{ID: "abc", Name: "web", State: "running"},
			},
		},
		"empty name": {
			given: given{
				input: InspectContainerInput{Name: ""},
				err:   errors.New("validation error"),
			},
			want: want{},
		},
		"service error": {
			given: given{
				input: InspectContainerInput{Name: "web"},
				err:   errors.New("docker error"),
			},
			want: want{
				called: true,
				name:   "web",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := containermock.NewMockContainerService(t)

			if test.want.called {
				mockService.On("InspectContainer", mock.Anything, container.InspectContainerParams{
					Name: test.want.name,
				}).Return(test.given.result, test.given.err)
			}

			handler := NewToolsHandler(mockService)

			_, output, err := handler.InspectContainer(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want.container, output.Container)
		})
	}
}
