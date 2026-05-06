package image

import (
	"context"
	"errors"
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
		"error": {
			given: given{
				params: providers.RemoveImageParams{Ref: "nginx"},
				err:    errors.New("docker error"),
			},
			want: want{
				ref: "nginx",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ImageRemove", mock.Anything, test.want.ref, client.ImageRemoveOptions{
				Force:         test.want.force,
				PruneChildren: test.want.pruneChildren,
				Platforms:     test.want.platforms,
			}).Return(client.ImageRemoveResult{}, test.given.err)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.RemoveImage(context.Background(), test.given.params)

			if test.given.err != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
