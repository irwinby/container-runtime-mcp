package image

import (
	"context"
	"errors"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	testdocker "github.com/irwinby/container-runtime-mcp/internal/testing/docker"
	"github.com/moby/moby/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderPullImage(t *testing.T) {
	type given struct {
		params   providers.PullImageParams
		pullErr  error
		waitErr  error
		closeErr error
	}

	type want struct {
		ref       string
		all       bool
		platforms []ocispec.Platform
		err       bool
	}

	plat := &ocispec.Platform{OS: "linux", Architecture: "amd64"}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.PullImageParams{
					Ref:      "nginx:latest",
					All:      true,
					Platform: plat,
				},
			},
			want: want{
				ref:       "nginx:latest",
				all:       true,
				platforms: []ocispec.Platform{*plat},
			},
		},
		"pull error": {
			given: given{
				params: providers.PullImageParams{
					Ref: "nginx:latest",
				},
				pullErr: errors.New("pull failed"),
			},
			want: want{
				ref: "nginx:latest",
				err: true,
			},
		},
		"wait error": {
			given: given{
				params: providers.PullImageParams{
					Ref: "nginx:latest",
				},
				waitErr: errors.New("wait failed"),
			},
			want: want{
				ref: "nginx:latest",
				err: true,
			},
		},
		"close error": {
			given: given{
				params: providers.PullImageParams{
					Ref: "nginx:latest",
				},
				closeErr: errors.New("close failed"),
			},
			want: want{
				ref: "nginx:latest",
				err: true,
			},
		},
		"success without platform": {
			given: given{
				params: providers.PullImageParams{
					Ref: "nginx:latest",
					All: true,
				},
			},
			want: want{
				ref: "nginx:latest",
				all: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			var response *testdocker.ProgressResponse

			if test.given.pullErr == nil {
				response = &testdocker.ProgressResponse{WaitErr: test.given.waitErr, CloseErr: test.given.closeErr}
			}

			mockClient.On("ImagePull", mock.Anything, test.want.ref, client.ImagePullOptions{
				All:       test.want.all,
				Platforms: test.want.platforms,
			}).Return(response, test.given.pullErr)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.PullImage(context.Background(), test.given.params)

			if test.want.err {
				require.Error(t, err)
				if response != nil {
					require.True(t, response.Closed, "response should be closed")
				}

				return
			}

			require.NoError(t, err)
			require.NotNil(t, response)
			require.True(t, response.Closed, "response should be closed")
		})
	}
}
