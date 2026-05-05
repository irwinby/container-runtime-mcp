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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := containermock.NewMockContainerService(t)
			if tt.want.called {
				mockSvc.On("InspectContainer", mock.Anything, container.InspectContainerParams{
					Name: tt.want.name,
				}).Return(tt.given.result, tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, output, err := handler.InspectContainer(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.container, output.Container)
		})
	}
}
