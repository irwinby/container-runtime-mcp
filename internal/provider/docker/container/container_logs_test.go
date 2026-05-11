package container

import (
	"bytes"
	"context"
	"errors"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	testdocker "github.com/irwinby/container-runtime-mcp/internal/testing/docker"
	testio "github.com/irwinby/container-runtime-mcp/internal/testing/io"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProviderContainerLogs(t *testing.T) {
	type given struct {
		params    providers.ContainerLogsParams
		clientErr error
		logData   []byte
	}

	type want struct {
		stdout string
		stderr string
		err    bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success": {
			given: given{
				params: providers.ContainerLogsParams{
					Name:   "web",
					Stdout: true,
					Stderr: true,
				},
				logData: func() []byte {
					var buf bytes.Buffer
					testdocker.WriteFrame(&buf, 1, []byte("hello"))
					testdocker.WriteFrame(&buf, 2, []byte("warn"))
					return buf.Bytes()
				}(),
			},
			want: want{stdout: "hello", stderr: "warn"},
		},
		"client error": {
			given: given{
				params: providers.ContainerLogsParams{
					Name:   "web",
					Stdout: true,
					Stderr: true,
				},
				clientErr: errors.New("docker error"),
			},
			want: want{err: true},
		},
		"stdout only": {
			given: given{
				params: providers.ContainerLogsParams{
					Name:   "web",
					Stdout: true,
					Stderr: false,
				},
				logData: func() []byte {
					var buf bytes.Buffer
					testdocker.WriteFrame(&buf, 1, []byte("hello"))
					return buf.Bytes()
				}(),
			},
			want: want{stdout: "hello"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			var response *testio.ReadCloser

			if test.given.clientErr == nil {
				response = &testio.ReadCloser{Data: test.given.logData}
			}

			mockClient.On("ContainerLogs", mock.Anything, test.given.params.Name, client.ContainerLogsOptions{
				ShowStdout: test.given.params.Stdout,
				ShowStderr: test.given.params.Stderr,
				Since:      test.given.params.Since,
				Timestamps: test.given.params.Timestamps,
				Tail:       test.given.params.Tail,
				Follow:     false,
			}).Return(response, test.given.clientErr)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.ContainerLogs(context.Background(), test.given.params)

			if test.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want.stdout, result.Stdout)
			require.Equal(t, test.want.stderr, result.Stderr)
			require.NotNil(t, response)
			require.True(t, response.Closed)
		})
	}
}
