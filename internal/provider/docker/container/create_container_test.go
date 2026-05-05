package container

import (
	"context"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderCreateContainer(t *testing.T) {
	type given struct {
		params providers.CreateContainerParams
		err    error
	}

	type want struct {
		id    string
		name  string
		image string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.CreateContainerParams{
					Name:  "web",
					Image: "nginx:latest",
				},
			},
			want: want{
				id:    "abc123",
				name:  "web",
				image: "nginx:latest",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ContainerCreate", mock.Anything, client.ContainerCreateOptions{
				Name: tt.want.name,
				Config: &container.Config{
					Image: tt.want.image,
				},
			}).Return(client.ContainerCreateResult{ID: tt.want.id}, tt.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			id, err := provider.CreateContainer(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.id, id)
		})
	}
}
