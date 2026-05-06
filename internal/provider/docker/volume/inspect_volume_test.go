package volume

import (
	"context"
	"errors"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/api/types/volume"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderInspectVolume(t *testing.T) {
	type given struct {
		params providers.InspectVolumeParams
		err    error
	}

	type want struct {
		name string
		vol  providers.VolumeInspect
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{params: providers.InspectVolumeParams{Name: "vol1"}},
			want:  want{name: "vol1", vol: providers.NewVolumeInspect().SetName("vol1").SetDriver("local")},
		},
		"error": {
			given: given{params: providers.InspectVolumeParams{Name: "vol1"}, err: errors.New("docker error")},
			want:  want{name: "vol1"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)
			mockClient.On("VolumeInspect", mock.Anything, test.want.name, client.VolumeInspectOptions{}).
				Return(client.VolumeInspectResult{Volume: volume.Volume{Name: "vol1", Driver: "local"}}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.InspectVolume(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want.vol.Name, result.Name)
			require.Equal(t, test.want.vol.Driver, result.Driver)
		})
	}
}
