package volume

import (
	"context"
	"errors"
	"testing"

	volumemock "github.com/irwinby/container-runtime-mcp/internal/api/mcp/handler/volume/mock"
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
		"dangling false": {
			given: given{
				input:  ListVolumesInput{Dangling: false},
				result: []volume.Volume{{Name: "vol1", Driver: "local"}},
			},
			want: want{
				called:   true,
				dangling: false,
				volumes:  []volume.Volume{{Name: "vol1", Driver: "local"}},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := volumemock.NewMockVolumeService(t)

			if test.want.called {
				mockService.On("ListVolumes", mock.Anything, volume.ListVolumesParams{
					Dangling: test.want.dangling,
				}).Return(test.given.result, test.given.err)
			}

			handler := NewToolsHandler(mockService)

			_, output, err := handler.ListVolumes(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want.volumes, output.Volumes)
		})
	}
}
