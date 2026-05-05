package container

import (
	"context"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderRemoveContainer(t *testing.T) {
	type given struct {
		params providers.RemoveContainerParams
		err    error
	}

	type want struct {
		name          string
		force         bool
		removeVolumes bool
		removeLinks   bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.RemoveContainerParams{
					Name:          "web",
					Force:         true,
					RemoveVolumes: true,
					RemoveLinks:   true,
				},
			},
			want: want{
				name:          "web",
				force:         true,
				removeVolumes: true,
				removeLinks:   true,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ContainerRemove", mock.Anything, tt.want.name, client.ContainerRemoveOptions{
				Force:         tt.want.force,
				RemoveVolumes: tt.want.removeVolumes,
				RemoveLinks:   tt.want.removeLinks,
			}).Return(client.ContainerRemoveResult{}, tt.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.RemoveContainer(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
