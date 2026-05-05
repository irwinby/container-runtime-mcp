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
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ContainerList", mock.Anything, client.ContainerListOptions{
				All:    tt.want.all,
				Limit:  tt.want.limit,
				Size:   tt.want.size,
				Latest: tt.want.latest,
			}).Return(client.ContainerListResult{
				Items: []container.Summary{{ID: "c1", Names: []string{"web"}}},
			}, tt.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.ListContainers(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, result, 1)
			require.Equal(t, tt.want.containers[0].ID, result[0].ID)
		})
	}
}
