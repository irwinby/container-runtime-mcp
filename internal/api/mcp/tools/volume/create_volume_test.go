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

func TestHandlerCreateVolume(t *testing.T) {
	type given struct {
		input  CreateVolumeInput
		result volume.VolumeInspect
		err    error
	}

	type want struct {
		called bool
		name   string
		driver string
		vol    volume.VolumeInspect
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"valid input": {
			given: given{
				input:  CreateVolumeInput{Name: "vol1", Driver: "local"},
				result: volume.VolumeInspect{Name: "vol1", Driver: "local"},
			},
			want: want{
				called: true,
				name:   "vol1",
				driver: "local",
				vol:    volume.VolumeInspect{Name: "vol1", Driver: "local"},
			},
		},
		"service error": {
			given: given{
				input: CreateVolumeInput{Name: "vol1"},
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
				mockSvc.On("CreateVolume", mock.Anything, volume.CreateVolumeParams{
					Name:   tt.want.name,
					Driver: tt.want.driver,
				}).Return(tt.given.result, tt.given.err)
			}

			handler := NewToolsHandler(mockSvc)

			_, output, err := handler.CreateVolume(context.Background(), &mcp.CallToolRequest{}, tt.given.input)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.vol, output.Volume)
		})
	}
}
