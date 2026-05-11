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

func TestProviderPushImage(t *testing.T) {
	type given struct {
		params   providers.PushImageParams
		pushErr  error
		waitErr  error
		closeErr error
	}

	type want struct {
		ref      string
		all      bool
		platform *ocispec.Platform
		err      bool
	}

	plat := &ocispec.Platform{OS: "linux", Architecture: "arm64"}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.PushImageParams{
					Ref:      "myapp:latest",
					All:      true,
					Platform: plat,
				},
			},
			want: want{
				ref:      "myapp:latest",
				all:      true,
				platform: plat,
			},
		},
		"push error": {
			given: given{
				params: providers.PushImageParams{
					Ref: "myapp:latest",
				},
				pushErr: errors.New("push failed"),
			},
			want: want{
				ref: "myapp:latest",
				err: true,
			},
		},
		"wait error": {
			given: given{
				params: providers.PushImageParams{
					Ref: "myapp:latest",
				},
				waitErr: errors.New("wait failed"),
			},
			want: want{
				ref: "myapp:latest",
				err: true,
			},
		},
		"close error": {
			given: given{
				params: providers.PushImageParams{
					Ref: "myapp:latest",
				},
				closeErr: errors.New("close failed"),
			},
			want: want{
				ref: "myapp:latest",
				err: true,
			},
		},
		"success without platform": {
			given: given{
				params: providers.PushImageParams{
					Ref: "myapp:latest",
					All: true,
				},
			},
			want: want{
				ref: "myapp:latest",
				all: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			var response *testdocker.ProgressResponse

			if test.given.pushErr == nil {
				response = &testdocker.ProgressResponse{WaitErr: test.given.waitErr, CloseErr: test.given.closeErr}
			}

			mockClient.On("ImagePush", mock.Anything, test.want.ref, client.ImagePushOptions{
				All:      test.want.all,
				Platform: test.want.platform,
			}).Return(response, test.given.pushErr)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.PushImage(context.Background(), test.given.params)

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
