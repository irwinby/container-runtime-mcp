// Package net provides test helpers for working with network resources.
package net

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// FreeTCPAddr returns an available TCP address on 127.0.0.1.
// The listener is closed immediately so the port can be reused.
func FreeTCPAddr(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	addr := listener.Addr().String()
	require.NoError(t, listener.Close())

	return addr
}

// RequireTCPListening waits until the given address is accepting TCP
// connections or the test timeout is reached.
func RequireTCPListening(t *testing.T, addr string) {
	t.Helper()

	require.Eventually(t, func() bool {
		connection, err := net.DialTimeout("tcp", addr, 50*time.Millisecond)
		if err != nil {
			return false
		}

		return connection.Close() == nil
	}, 2*time.Second, 10*time.Millisecond)
}
