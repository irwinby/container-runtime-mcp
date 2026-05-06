package image

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

func TestProviderTagImage(t *testing.T) {
	type given struct {
		params providers.TagImageParams
		err    error
	}

	type want struct {
		source string
		target string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.TagImageParams{
					Source: "nginx:latest",
					Target: "my-nginx:latest",
				},
			},
			want: want{
				source: "nginx:latest",
				target: "my-nginx:latest",
			},
		},
		"error": {
			given: given{
				params: providers.TagImageParams{
					Source: "nginx:latest",
					Target: "my-nginx:latest",
				},
				err: errors.New("docker error"),
			},
			want: want{
				source: "nginx:latest",
				target: "my-nginx:latest",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ImageTag", mock.Anything, client.ImageTagOptions{
				Source: test.want.source,
				Target: test.want.target,
			}).Return(client.ImageTagResult{}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.TagImage(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
