// Package iotest provides generic I/O test doubles.
package iotest

import "io"

// ReadCloser is a test double for io.ReadCloser backed by a byte slice.
type ReadCloser struct {
	Data   []byte
	Closed bool
}

// Read implements io.Reader.
func (r *ReadCloser) Read(p []byte) (int, error) {
	if len(r.Data) == 0 {
		return 0, io.EOF
	}

	n := copy(p, r.Data)
	r.Data = r.Data[n:]

	return n, nil
}

// Close implements io.Closer.
func (r *ReadCloser) Close() error {
	r.Closed = true
	return nil
}
