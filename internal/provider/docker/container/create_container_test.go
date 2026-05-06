package container

import (
	"context"
	"errors"
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
		"error": {
			given: given{
				params: providers.CreateContainerParams{
					Name:  "web",
					Image: "nginx:latest",
				},
				err: errors.New("docker error"),
			},
			want: want{
				name:  "web",
				image: "nginx:latest",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ContainerCreate", mock.Anything, client.ContainerCreateOptions{
				Name: test.want.name,
				Config: &container.Config{
					Image: test.want.image,
				},
			}).Return(client.ContainerCreateResult{ID: test.want.id}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			id, err := provider.CreateContainer(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want.id, id)
		})
	}
}
