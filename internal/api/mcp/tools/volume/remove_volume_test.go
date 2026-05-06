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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := volumemock.NewMockVolumeService(t)

			if tt.want.called {
				mockSvc.On("RemoveVolume", mock.Anything, volume.RemoveVolumeParams{
					Name:  tt.want.name,
					Force: tt.want.force,
				}).Return(tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, _, err := handler.RemoveVolume(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			if !tt.want.called {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
