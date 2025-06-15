package abstractions

import (
	u "net/url"
)

// Request represents an abstract HTTP request.
type Request struct {
	// The HTTP method of the request.
	Method HttpMethod
	uri    *u.URL
	// The Request Headers.
	Headers *RequestHeaders
	// The Query Parameters of the request.
	QueryParameters map[string]any
	// The Request Body.
	Content []byte
	// The path parameters to use for the URL template when generating the URI.
	PathParameters map[string]any
	// The Url template for the current request.
	UrlTemplate string
	options     map[string]RequestOption
}

const raw_url_key = "request-raw-url"

// NewRequestInformation creates a new RequestInformation object with default values.
func NewRequestInformation() *Request {
	return &Request{
		Headers:         NewRequestHeaders(),
		QueryParameters: make(map[string]any),
		options:         make(map[string]RequestOption),
		PathParameters:  make(map[string]any),
	}
}

func castItem[T any, R interface{}](collection []T, mutator func(t T) R) []R {
	if len(collection) > 0 {
		cast := make([]R, len(collection))
		for i, v := range collection {
			cast[i] = mutator(v)
		}
		return cast
	}
	return nil
}

const contentTypeHeader = "Content-Type"
const binaryContentType = "application/octet-stream"
const observabilityTracerName = "github.com/microsoft/kiota-abstractions-go"
