package util

import (
	"io"
	"net/http"
)

// ReadCloser is io.ReadCloser, which is just for generating mock.
type ReadCloser interface{ io.ReadCloser }

// HTTPClient is an interface of net/http's Client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
