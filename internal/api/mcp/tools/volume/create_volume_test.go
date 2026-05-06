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
		"with driver opts and labels": {
			given: given{
				input: CreateVolumeInput{
					Name:       "vol1",
					Driver:     "local",
					DriverOpts: map[string]string{"size": "10Gi"},
					Labels:     map[string]string{"env": "test"},
				},
				result: volume.VolumeInspect{Name: "vol1", Driver: "local"},
			},
			want: want{
				called: true,
				name:   "vol1",
				driver: "local",
				vol:    volume.VolumeInspect{Name: "vol1", Driver: "local"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := volumemock.NewMockVolumeService(t)

			if test.want.called {
				mockService.On("CreateVolume", mock.Anything, volume.CreateVolumeParams{
					Name:       test.want.name,
					Driver:     test.want.driver,
					DriverOpts: test.given.input.DriverOpts,
					Labels:     test.given.input.Labels,
				}).Return(test.given.result, test.given.err)
			}

			handler := NewToolsHandler(mockService)

			_, output, err := handler.CreateVolume(context.Background(), &mcp.CallToolRequest{}, test.given.input)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want.vol, output.Volume)
		})
	}
}
