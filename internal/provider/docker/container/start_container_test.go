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

func TestProviderStartContainer(t *testing.T) {
	type given struct {
		params providers.StartContainerParams
		err    error
	}

	type want struct {
		name string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.StartContainerParams{Name: "web"},
			},
			want: want{
				name: "web",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ContainerStart", mock.Anything, tt.want.name, client.ContainerStartOptions{}).Return(client.ContainerStartResult{}, tt.given.err)

			provider := NewProvider(mockClient, nopTimeout)
			err := provider.StartContainer(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
