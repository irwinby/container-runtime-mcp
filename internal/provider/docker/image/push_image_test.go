package image

import (
	"context"
	"errors"
	"io"
	"iter"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/api/types/jsonstream"
	"github.com/moby/moby/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testPushResponse struct {
	waitErr  error
	closeErr error
	closed   bool
}

func (r *testPushResponse) Read(p []byte) (int, error) { return 0, io.EOF }
func (r *testPushResponse) Close() error {
	r.closed = true
	return r.closeErr
}
func (r *testPushResponse) JSONMessages(ctx context.Context) iter.Seq2[jsonstream.Message, error] {
	return func(yield func(jsonstream.Message, error) bool) {}
}
func (r *testPushResponse) Wait(ctx context.Context) error { return r.waitErr }

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

			var resp *testPushResponse
			if test.given.pushErr == nil {
				resp = &testPushResponse{waitErr: test.given.waitErr, closeErr: test.given.closeErr}
			}

			mockClient.On("ImagePush", mock.Anything, test.want.ref, client.ImagePushOptions{
				All:      test.want.all,
				Platform: test.want.platform,
			}).Return(resp, test.given.pushErr)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.PushImage(context.Background(), test.given.params)

			if test.given.pushErr != nil || test.given.waitErr != nil || test.given.closeErr != nil {
				require.Error(t, err)
				if resp != nil {
					require.True(t, resp.closed, "response should be closed")
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.True(t, resp.closed, "response should be closed")
		})
	}
}
