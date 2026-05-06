package container

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func writeDockerLogFrame(buf *bytes.Buffer, streamType byte, data []byte) {
	var header [8]byte
	header[0] = streamType
	binary.BigEndian.PutUint32(header[4:], uint32(len(data)))
	buf.Write(header[:])
	buf.Write(data)
}

type testReadCloser struct {
	data   []byte
	closed bool
}

func (r *testReadCloser) Read(p []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.data)
	r.data = r.data[n:]
	return n, nil
}

func (r *testReadCloser) Close() error {
	r.closed = true
	return nil
}

func TestProviderContainerLogs(t *testing.T) {
	tests := map[string]struct {
		params       providers.ContainerLogsParams
		clientErr    error
		logData      []byte
		wantStdout   string
		wantStderr   string
		wantErr      bool
		wantStderrOk bool
	}{
		"success": {
			params: providers.ContainerLogsParams{
				Name:   "web",
				Stdout: true,
				Stderr: true,
			},
			logData: func() []byte {
				var buf bytes.Buffer
				writeDockerLogFrame(&buf, 1, []byte("hello"))
				writeDockerLogFrame(&buf, 2, []byte("warn"))
				return buf.Bytes()
			}(),
			wantStdout: "hello",
			wantStderr: "warn",
		},
		"client error": {
			params: providers.ContainerLogsParams{
				Name:   "web",
				Stdout: true,
				Stderr: true,
			},
			clientErr: errors.New("docker error"),
			wantErr:   true,
		},
		"stdout only": {
			params: providers.ContainerLogsParams{
				Name:   "web",
				Stdout: true,
				Stderr: false,
			},
			logData: func() []byte {
				var buf bytes.Buffer
				writeDockerLogFrame(&buf, 1, []byte("hello"))
				return buf.Bytes()
			}(),
			wantStdout: "hello",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			var resp *testReadCloser
			if test.clientErr == nil {
				resp = &testReadCloser{data: test.logData}
			}

			mockClient.On("ContainerLogs", mock.Anything, test.params.Name, client.ContainerLogsOptions{
				ShowStdout: test.params.Stdout,
				ShowStderr: test.params.Stderr,
				Since:      test.params.Since,
				Timestamps: test.params.Timestamps,
				Tail:       test.params.Tail,
				Follow:     false,
			}).Return(resp, test.clientErr)

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.ContainerLogs(context.Background(), test.params)

			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.wantStdout, result.Stdout)
			require.Equal(t, test.wantStderr, result.Stderr)
			require.NotNil(t, resp)
			require.True(t, resp.closed)
		})
	}
}
