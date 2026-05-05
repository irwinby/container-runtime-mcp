package image

import (
	"context"
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
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ImageInspect", mock.Anything, tt.want.ref).Return(client.ImageInspectResult{
				InspectResponse: image.InspectResponse{
					ID:           tt.want.id,
					RepoTags:     []string{"nginx:latest"},
					Size:         1024,
					Created:      "2024-01-01",
					Architecture: tt.want.architecture,
					Os:           tt.want.os,
				},
			}, tt.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.InspectImage(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.id, result.ID)
			require.Equal(t, tt.want.architecture, result.Architecture)
			require.Equal(t, tt.want.os, result.OS)
		})
	}
}
