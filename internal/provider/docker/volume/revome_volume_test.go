package volume

import (
	"context"
	"errors"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderRemoveVolume(t *testing.T) {
	type given struct {
		params providers.RemoveVolumeParams
		err    error
	}

	type want struct {
		name  string
		force bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{params: providers.RemoveVolumeParams{Name: "vol1", Force: true}},
			want:  want{name: "vol1", force: true},
		},
		"error": {
			given: given{params: providers.RemoveVolumeParams{Name: "vol1"}, err: errors.New("docker error")},
			want:  want{name: "vol1"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("VolumeRemove", mock.Anything, test.want.name, client.VolumeRemoveOptions{
				Force: test.want.force,
			}).Return(client.VolumeRemoveResult{}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.RemoveVolume(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
