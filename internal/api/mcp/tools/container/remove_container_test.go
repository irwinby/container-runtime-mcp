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

func TestHandlerRemoveContainer(t *testing.T) {
	type given struct {
		input RemoveContainerInput
		err   error
	}

	type want struct {
		called        bool
		name          string
		force         bool
		removeVolumes bool
		removeLinks   bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				input: RemoveContainerInput{Name: "web"},
			},
			want: want{
				called: true,
				name:   "web",
			},
		},
		"whitespace name": {
			given: given{
				input: RemoveContainerInput{Name: "  web  "},
			},
			want: want{
				called: true,
				name:   "web",
			},
		},
		"with options": {
			given: given{
				input: RemoveContainerInput{Name: "web", Force: true, RemoveVolumes: true, RemoveLinks: true},
			},
			want: want{
				called:        true,
				name:          "web",
				force:         true,
				removeVolumes: true,
				removeLinks:   true,
			},
		},
		"empty name": {
			given: given{
				input: RemoveContainerInput{Name: ""},
				err:   errors.New("validation error"),
			},
			want: want{},
		},
		"service error": {
			given: given{
				input: RemoveContainerInput{Name: "web"},
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
				mockSvc.On("RemoveContainer", mock.Anything, container.RemoveContainerParams{
					Name:          tt.want.name,
					Force:         tt.want.force,
					RemoveVolumes: tt.want.removeVolumes,
					RemoveLinks:   tt.want.removeLinks,
				}).Return(tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, _, err := handler.RemoveContainer(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
