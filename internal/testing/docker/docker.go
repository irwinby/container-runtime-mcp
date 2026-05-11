// Package docker provides test helpers for Docker SDK interactions.
package docker

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"iter"
	"net"
	"testing"

	"github.com/moby/moby/api/types/jsonstream"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/require"
)

// WriteFrame writes a Docker multiplexed stream frame into buf.
// streamType: 1 = stdout, 2 = stderr.
func WriteFrame(buf *bytes.Buffer, streamType byte, data []byte) {
	var header [8]byte
	header[0] = streamType
	binary.BigEndian.PutUint32(header[4:], uint32(len(data)))
	buf.Write(header[:])
	buf.Write(data)
}

// SetupHijackedConn creates a hijacked connection pair for testing
// container exec attach. When attachStdin is false, a net.Pipe is used.
// When true, a TCP loopback connection is used.
func SetupHijackedConn(t *testing.T, attachStdin bool) (serverConnection net.Conn, response client.HijackedResponse) {
	t.Helper()

	if !attachStdin {
		serverConnection, clientConnection := net.Pipe()

		t.Cleanup(func() {
			_ = serverConnection.Close()
			_ = clientConnection.Close()
		})

		return serverConnection, client.HijackedResponse{
			Conn:   clientConnection,
			Reader: bufio.NewReader(clientConnection),
		}
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	t.Cleanup(func() { _ = listener.Close() })

	connections := make(chan net.Conn, 1)

	go func() {
		connection, err := listener.Accept()
		if err != nil {
			close(connections)
			return
		}

		connections <- connection
	}()

	clientConnection, err := net.Dial("tcp", listener.Addr().String())
	require.NoError(t, err)

	serverConnection = <-connections
	require.NotNil(t, serverConnection)

	t.Cleanup(func() {
		_ = serverConnection.Close()
		_ = clientConnection.Close()
	})

	return serverConnection, client.HijackedResponse{
		Conn:   clientConnection,
		Reader: bufio.NewReader(clientConnection),
	}
}

// ProgressResponse is a test double for Docker image pull/push responses.
type ProgressResponse struct {
	WaitErr  error
	CloseErr error
	Closed   bool
}

// Read implements io.Reader.
func (r *ProgressResponse) Read(_ []byte) (int, error) { return 0, io.EOF }

// Close implements io.Closer.
func (r *ProgressResponse) Close() error {
	r.Closed = true
	return r.CloseErr
}

// JSONMessages implements the JSON message iterator for pull/push responses.
func (r *ProgressResponse) JSONMessages(_ context.Context) iter.Seq2[jsonstream.Message, error] {
	return func(yield func(jsonstream.Message, error) bool) {}
}

// Wait blocks until the operation completes.
func (r *ProgressResponse) Wait(_ context.Context) error { return r.WaitErr }
