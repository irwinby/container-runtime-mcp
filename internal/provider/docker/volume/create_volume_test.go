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

func TestProviderCreateVolume(t *testing.T) {
	type given struct {
		params providers.CreateVolumeParams
		err    error
	}

	type want struct {
		vol providers.VolumeInspect
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{params: providers.CreateVolumeParams{Name: "vol1", Driver: "local"}},
			want:  want{vol: providers.NewVolumeInspect().SetName("vol1").SetDriver("local")},
		},
		"error": {
			given: given{params: providers.CreateVolumeParams{Name: "vol1"}, err: errors.New("docker error")},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("VolumeCreate", mock.Anything, client.VolumeCreateOptions{
				Name:   test.given.params.Name,
				Driver: test.given.params.Driver,
			}).Return(client.VolumeCreateResult{Volume: volume.Volume{Name: test.given.params.Name, Driver: test.given.params.Driver}}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.CreateVolume(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want.vol.Name, result.Name)
		})
	}
}
