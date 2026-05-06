package container

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"net"
	"testing"

	providers "github.com/irwinby/container-runtime-mcp/internal/provider"
	dockermock "github.com/irwinby/container-runtime-mcp/internal/provider/docker/mock"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// writeStdFrame writes a Docker multiplexed stream frame.
// streamType: 1 = stdout, 2 = stderr.
func writeStdFrame(buf *bytes.Buffer, streamType byte, data []byte) {
	var header [8]byte
	header[0] = streamType
	binary.BigEndian.PutUint32(header[4:], uint32(len(data)))
	buf.Write(header[:])
	buf.Write(data)
}

func setupHijackedConn(t *testing.T, attachStdin bool) (server net.Conn, resp client.HijackedResponse) {
	t.Helper()

	if !attachStdin {
		server, clientConn := net.Pipe()
		t.Cleanup(func() { server.Close(); clientConn.Close() })
		return server, client.HijackedResponse{
			Conn:   clientConn,
			Reader: bufio.NewReader(clientConn),
		}
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { listener.Close() })

	serverCh := make(chan net.Conn, 1)
	go func() {
		s, err := listener.Accept()
		if err != nil {
			close(serverCh)
			return
		}
		serverCh <- s
	}()

	clientConn, err := net.Dial("tcp", listener.Addr().String())
	require.NoError(t, err)

	server = <-serverCh
	require.NotNil(t, server)
	t.Cleanup(func() {
		server.Close()
		clientConn.Close()
	})

	return server, client.HijackedResponse{
		Conn:   clientConn,
		Reader: bufio.NewReader(clientConn),
	}
}

func TestProviderExecContainer(t *testing.T) {
	type given struct {
		params        providers.ExecContainerParams
		createErr     error
		attachErr     error
		inspectErr    error
		attachOutput  []byte
		inspectResult client.ExecInspectResult
	}

	type want struct {
		exitCode int
		stdout   string
		stderr   string
		err      bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"success non-tty": {
			given: given{
				params: providers.ExecContainerParams{
					Name:         "web",
					Cmd:          []string{"echo", "hello"},
					AttachStdout: true,
					AttachStderr: true,
				},
				attachOutput: func() []byte {
					var buf bytes.Buffer
					writeStdFrame(&buf, 1, []byte("hello"))
					writeStdFrame(&buf, 2, []byte("warn"))
					return buf.Bytes()
				}(),
				inspectResult: client.ExecInspectResult{ExitCode: 0},
			},
			want: want{exitCode: 0, stdout: "hello", stderr: "warn"},
		},
		"success tty": {
			given: given{
				params: providers.ExecContainerParams{
					Name:         "web",
					Cmd:          []string{"echo", "hello"},
					TTY:          true,
					AttachStdout: true,
				},
				attachOutput:  []byte("hello"),
				inspectResult: client.ExecInspectResult{ExitCode: 0},
			},
			want: want{exitCode: 0, stdout: "hello"},
		},
		"create error": {
			given: given{
				params:    providers.ExecContainerParams{Name: "web", Cmd: []string{"echo"}},
				createErr: errors.New("docker error"),
			},
			want: want{err: true},
		},
		"attach error": {
			given: given{
				params:    providers.ExecContainerParams{Name: "web", Cmd: []string{"echo"}},
				attachErr: errors.New("docker error"),
			},
			want: want{err: true},
		},
		"inspect error": {
			given: given{
				params:     providers.ExecContainerParams{Name: "web", Cmd: []string{"echo"}},
				inspectErr: errors.New("docker error"),
			},
			want: want{err: true},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockClient := dockermock.NewMockDockerClient(t)

			mockClient.On("ExecCreate", mock.Anything, test.given.params.Name, client.ExecCreateOptions{
				Cmd:          test.given.params.Cmd,
				Env:          test.given.params.Env,
				WorkingDir:   test.given.params.WorkingDir,
				User:         test.given.params.User,
				Privileged:   test.given.params.Privileged,
				TTY:          test.given.params.TTY,
				AttachStdin:  test.given.params.AttachStdin,
				AttachStdout: test.given.params.AttachStdout,
				AttachStderr: test.given.params.AttachStderr,
			}).Return(client.ExecCreateResult{ID: "exec-123"}, test.given.createErr)

			if test.given.createErr == nil {
				var attachResult client.ExecAttachResult
				if test.given.attachErr == nil {
					server, resp := setupHijackedConn(t, test.given.params.AttachStdin)

					go func() {
						_, _ = server.Write(test.given.attachOutput)
						_ = server.Close()
					}()

					attachResult = client.ExecAttachResult{
						HijackedResponse: resp,
					}
				}

				mockClient.On("ExecAttach", mock.Anything, "exec-123", client.ExecAttachOptions{
					TTY: test.given.params.TTY,
				}).Return(attachResult, test.given.attachErr)

				if test.given.attachErr == nil {
					mockClient.On("ExecInspect", mock.Anything, "exec-123", client.ExecInspectOptions{}).
						Return(test.given.inspectResult, test.given.inspectErr)
				}
			}

			provider := NewProvider(mockClient, nopTimeout)

			result, err := provider.ExecContainer(context.Background(), test.given.params)

			if test.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want.exitCode, result.ExitCode)
			require.Equal(t, test.want.stdout, result.Stdout)
			require.Equal(t, test.want.stderr, result.Stderr)
		})
	}
}
