package system

import (
	"context"
	"testing"

	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/api/types/system"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderSystemInfo(t *testing.T) {
	type given struct {
		err error
	}

	type want struct {
		id                string
		containers        int
		containersRunning int
		images            int
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{},
			want: want{
				id:                "daemon1",
				containers:        5,
				containersRunning: 3,
				images:            10,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("Info", mock.Anything, client.InfoOptions{}).Return(client.SystemInfoResult{
				Info: system.Info{
					ID:                tt.want.id,
					Containers:        tt.want.containers,
					ContainersRunning: tt.want.containersRunning,
					Images:            tt.want.images,
				},
			}, tt.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.SystemInfo(context.Background())

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.id, result.ID)
			require.Equal(t, tt.want.containers, result.Containers)
			require.Equal(t, tt.want.containersRunning, result.ContainersRunning)
			require.Equal(t, tt.want.images, result.Images)
		})
	}
}
