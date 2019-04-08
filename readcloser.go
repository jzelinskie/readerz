package readerz

import (
	"bufio"
	"compress/gzip"
	"io"
)

// BufferedReadCloser implements a bufio.Reader that when calling Close() will
// also close the parent Reader.
type BufferedReadCloser struct {
	r io.ReadCloser
	*bufio.Reader
}

var _ io.ReadCloser = &BufferedReadCloser{}

// Close closes the parent io.ReadCloser.
func (r *BufferedReadCloser) Close() error {
	return r.Close()
}

// NewBufferedReadCloser returns a new BufferedReadCloser.
func NewBufferedReadCloser(r io.ReadCloser) *BufferedReadCloser {
	return &BufferedReadCloser{
		r,
		bufio.NewReader(r),
	}
}

// GzipReadCloser implements a gzip.Reader.
type GzipReadCloser struct {
	r *BufferedReadCloser
	*gzip.Reader
}

var _ io.ReadCloser = &GzipReadCloser{}

// Close closes the gzip.Reader and the parent io.ReadCloser.
func (z *GzipReadCloser) Close() error {
	if err := z.Reader.Close(); err != nil {
		return err
	}

	return z.r.Close()
}

// NewGzipReadCloser returns a new io.ReadCloser that will be wrapped in
// a gzip.Reader if the first bytes read are the gzip header.
func NewGzipReadCloser(r io.ReadCloser) (io.ReadCloser, error) {
	br := NewBufferedReadCloser(r)
	header, err := br.Peek(2)
	if err != nil {
		return nil, err
	}

	if header[0] == 0x1f && header[1] == 0x8b {
		zr, err := gzip.NewReader(br)
		if err != nil {
			return nil, err
		}
		return &GzipReadCloser{br, zr}, nil
	}

	return br, nil
}
