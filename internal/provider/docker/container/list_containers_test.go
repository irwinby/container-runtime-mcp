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

func TestProviderListContainers(t *testing.T) {
	type given struct {
		params providers.ListContainersParams
		err    error
	}

	type want struct {
		all        bool
		limit      int
		size       bool
		latest     bool
		containers []providers.Container
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.ListContainersParams{
					All:    true,
					Limit:  10,
					Size:   true,
					Latest: true,
				},
			},
			want: want{
				all:    true,
				limit:  10,
				size:   true,
				latest: true,
				containers: []providers.Container{
					providers.NewContainer().SetID("c1").SetNames([]string{"web"}),
				},
			},
		},
		"error": {
			given: given{
				params: providers.ListContainersParams{All: true},
				err:    errors.New("docker error"),
			},
			want: want{
				all: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ContainerList", mock.Anything, client.ContainerListOptions{
				All:    test.want.all,
				Limit:  test.want.limit,
				Size:   test.want.size,
				Latest: test.want.latest,
			}).Return(client.ContainerListResult{
				Items: []container.Summary{{ID: "c1", Names: []string{"web"}}},
			}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.ListContainers(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, result, 1)
			require.Equal(t, test.want.containers[0].ID, result[0].ID)
		})
	}
}
