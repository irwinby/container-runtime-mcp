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

type testPullResponse struct {
	waitErr  error
	closeErr error
	closed   bool
}

func (r *testPullResponse) Read(p []byte) (int, error) { return 0, io.EOF }
func (r *testPullResponse) Close() error {
	r.closed = true
	return r.closeErr
}
func (r *testPullResponse) JSONMessages(ctx context.Context) iter.Seq2[jsonstream.Message, error] {
	return func(yield func(jsonstream.Message, error) bool) {}
}
func (r *testPullResponse) Wait(ctx context.Context) error { return r.waitErr }

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
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			var resp *testPullResponse
			if tt.given.pullErr == nil {
				resp = &testPullResponse{waitErr: tt.given.waitErr, closeErr: tt.given.closeErr}
			}

			mockClient.On("ImagePull", mock.Anything, tt.want.ref, client.ImagePullOptions{
				All:       tt.want.all,
				Platforms: tt.want.platforms,
			}).Return(resp, tt.given.pullErr)

			provider := NewProvider(mockClient, nopTimeout)

			err := provider.PullImage(context.Background(), tt.given.params)

			if tt.given.pullErr != nil || tt.given.waitErr != nil || tt.given.closeErr != nil {
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
