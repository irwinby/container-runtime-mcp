package image

import (
	"context"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderRemoveImage(t *testing.T) {
	type given struct {
		params providers.RemoveImageParams
		err    error
	}

	type want struct {
		ref           string
		force         bool
		pruneChildren bool
		platforms     []ocispec.Platform
	}

	plat := &ocispec.Platform{OS: "linux", Architecture: "amd64"}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.RemoveImageParams{
					Ref:           "nginx",
					Force:         true,
					PruneChildren: true,
					Platform:      plat,
				},
			},
			want: want{
				ref:           "nginx",
				force:         true,
				pruneChildren: true,
				platforms:     []ocispec.Platform{*plat},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ImageRemove", mock.Anything, tt.want.ref, client.ImageRemoveOptions{
				Force:         tt.want.force,
				PruneChildren: tt.want.pruneChildren,
				Platforms:     tt.want.platforms,
			}).Return(client.ImageRemoveResult{}, tt.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.RemoveImage(context.Background(), tt.given.params)

			if tt.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
