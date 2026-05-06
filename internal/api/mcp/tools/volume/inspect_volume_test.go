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

func TestHandlerInspectVolume(t *testing.T) {
	type given struct {
		input  InspectVolumeInput
		result volume.VolumeInspect
		err    error
	}

	type want struct {
		called bool
		name   string
		vol    volume.VolumeInspect
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				input:  InspectVolumeInput{Name: "vol1"},
				result: volume.VolumeInspect{Name: "vol1", Driver: "local"},
			},
			want: want{
				called: true,
				name:   "vol1",
				vol:    volume.VolumeInspect{Name: "vol1", Driver: "local"},
			},
		},
		"empty name": {
			given: given{
				input: InspectVolumeInput{Name: ""},
				err:   errors.New("validation error"),
			},
			want: want{},
		},
		"service error": {
			given: given{
				input: InspectVolumeInput{Name: "vol1"},
				err:   errors.New("docker error"),
			},
			want: want{
				called: true,
				name:   "vol1",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockSvc := volumemock.NewMockVolumeService(t)

			if tt.want.called {
				mockSvc.On("InspectVolume", mock.Anything, volume.InspectVolumeParams{
					Name: tt.want.name,
				}).Return(tt.given.result, tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, output, err := handler.InspectVolume(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			if !tt.want.called {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.vol, output.Volume)
		})
	}
}
