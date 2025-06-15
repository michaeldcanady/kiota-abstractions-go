package abstractions

import "io"

// nopCloser is an alternate io.nopCloser implementation which
// provides io.ReadSeekCloser instead of io.ReadCloser as we need
// Seek for retries
type nopCloser struct {
	io.ReadSeeker
}

func NopCloser(r io.ReadSeeker) io.ReadSeekCloser {
	return nopCloser{r}
}

func (nopCloser) Close() error { return nil }
