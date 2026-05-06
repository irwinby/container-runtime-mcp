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

func TestProviderInspectImage(t *testing.T) {
	type given struct {
		params providers.InspectImageParams
		err    error
	}

	type want struct {
		ref          string
		id           string
		architecture string
		os           string
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.InspectImageParams{Ref: "nginx:latest"},
			},
			want: want{
				ref:          "nginx:latest",
				id:           "img1",
				architecture: "amd64",
				os:           "linux",
			},
		},
		"error": {
			given: given{
				params: providers.InspectImageParams{Ref: "nginx:latest"},
				err:    errors.New("docker error"),
			},
			want: want{
				ref: "nginx:latest",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ImageInspect", mock.Anything, test.want.ref).Return(client.ImageInspectResult{
				InspectResponse: image.InspectResponse{
					ID:           test.want.id,
					RepoTags:     []string{"nginx:latest"},
					Size:         1024,
					Created:      "2024-01-01",
					Architecture: test.want.architecture,
					Os:           test.want.os,
				},
			}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.InspectImage(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want.id, result.ID)
			require.Equal(t, test.want.architecture, result.Architecture)
			require.Equal(t, test.want.os, result.OS)
		})
	}
}
