package volume

import (
	"context"
	"errors"
	"testing"

	volumemock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/tools/volume/mock"
	"github.com/irwinby/container-runtime-mcp/internal/service/volume"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerListVolumes(t *testing.T) {
	type given struct {
		input  ListVolumesInput
		result []volume.Volume
		err    error
	}

	type want struct {
		called   bool
		dangling bool
		volumes  []volume.Volume
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				input:  ListVolumesInput{Dangling: true},
				result: []volume.Volume{{Name: "vol1", Driver: "local"}},
			},
			want: want{
				called:   true,
				dangling: true,
				volumes:  []volume.Volume{{Name: "vol1", Driver: "local"}},
			},
		},
		"service error": {
			given: given{
				input: ListVolumesInput{},
				err:   errors.New("docker error"),
			},
			want: want{called: true},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := volumemock.NewMockVolumeService(t)

			if tt.want.called {
				mockSvc.On("ListVolumes", mock.Anything, volume.ListVolumesParams{
					Dangling: tt.want.dangling,
				}).Return(tt.given.result, tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, output, err := handler.ListVolumes(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.volumes, output.Volumes)
		})
	}
}
