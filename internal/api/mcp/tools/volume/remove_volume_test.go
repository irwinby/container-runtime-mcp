package volume

import (
	"context"
	"errors"
	"testing"

	volumemock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/volume/mock"
	"github.com/irwinby/container-runtime-mcp/internal/service/volume"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerRemoveVolume(t *testing.T) {
	type given struct {
		input RemoveVolumeInput
		err   error
	}

	type want struct {
		called bool
		name   string
		force  bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{input: RemoveVolumeInput{Name: "vol1", Force: true}},
			want:  want{called: true, name: "vol1", force: true},
		},
		"empty name": {
			given: given{input: RemoveVolumeInput{Name: ""}, err: errors.New("validation error")},
			want:  want{},
		},
		"service error": {
			given: given{input: RemoveVolumeInput{Name: "vol1"}, err: errors.New("docker error")},
			want:  want{called: true, name: "vol1"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := volumemock.NewMockVolumeService(t)

			if test.want.called {
				mockService.On("RemoveVolume", mock.Anything, volume.RemoveVolumeParams{
					Name:  test.want.name,
					Force: test.want.force,
				}).Return(test.given.err)
			}

			handler := NewToolsHandler(mockService)

			_, _, err := handler.RemoveVolume(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			if !test.want.called {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
