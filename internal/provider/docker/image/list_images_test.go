package image

import (
	"context"
	"errors"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/api/types/image"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderListImages(t *testing.T) {
	type given struct {
		params providers.ListImagesParams
		err    error
	}

	type want struct {
		all        bool
		sharedSize bool
		images     []providers.Image
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.ListImagesParams{
					All:        true,
					SharedSize: true,
				},
			},
			want: want{
				all:        true,
				sharedSize: true,
				images: []providers.Image{
					providers.NewImage().SetID("img1").SetRepoTags([]string{"nginx:latest"}),
				},
			},
		},
		"error": {
			given: given{
				params: providers.ListImagesParams{All: true},
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

			mockClient.On("ImageList", mock.Anything, client.ImageListOptions{
				All:        test.want.all,
				SharedSize: test.want.sharedSize,
			}).Return(client.ImageListResult{
				Items: []image.Summary{{ID: "img1", RepoTags: []string{"nginx:latest"}}},
			}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.ListImages(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, result, 1)
			require.Equal(t, test.want.images[0].ID, result[0].ID)
		})
	}
}
