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
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			var resp *testPushResponse
			if tt.given.pushErr == nil {
				resp = &testPushResponse{waitErr: tt.given.waitErr, closeErr: tt.given.closeErr}
			}

			mockClient.On("ImagePush", mock.Anything, tt.want.ref, client.ImagePushOptions{
				All:      tt.want.all,
				Platform: tt.want.platform,
			}).Return(resp, tt.given.pushErr)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.PushImage(context.Background(), tt.given.params)

			if tt.given.pushErr != nil || tt.given.waitErr != nil || tt.given.closeErr != nil {
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
