package volume

import (
	"context"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/api/types/volume"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderListVolumes(t *testing.T) {
	type given struct {
		params providers.ListVolumesParams
		err    error
	}

	type want struct {
		dangling bool
		volumes  []providers.Volume
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{params: providers.ListVolumesParams{Dangling: true}},
			want: want{
				dangling: true,
				volumes: []providers.Volume{
					providers.NewVolume().SetName("vol1").SetDriver("local"),
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			filters := client.Filters{}
			if tt.want.dangling {
				filters = filters.Add("dangling", "true")
			}

			mockClient.On("VolumeList", mock.Anything, client.VolumeListOptions{
				Filters: filters,
			}).Return(client.VolumeListResult{
				Items: []volume.Volume{{Name: "vol1", Driver: "local"}},
			}, tt.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.ListVolumes(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, result, 1)
			require.Equal(t, tt.want.volumes[0].Name, result[0].Name)
			require.Equal(t, tt.want.volumes[0].Driver, result[0].Driver)
		})
	}
}
