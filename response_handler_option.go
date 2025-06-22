package abstractions

import nethttp "net/http"

// ResponseHandlerOption represents an abstract provider for ResponseHandler
type ResponseHandlerOption[T any] interface {
	GetResponseHandler() ResponseHandler[T]
	SetResponseHandler(responseHandler ResponseHandler[T])
	GetKey() ResponseOptionKey
}

var ResponseHandlerOptionKey = ResponseOptionKey{
	Key: "ResponseHandlerOptionKey",
}

type nativeResponseHandlerOption struct {
	handler ResponseHandler[nethttp.Response]
}

// NewRequestHandlerOption creates a new RequestInformation object with default values.
func NewRequestHandlerOption() ResponseHandlerOption[nethttp.Response] {
	return &nativeResponseHandlerOption{}
}

func (r *nativeResponseHandlerOption) GetResponseHandler() ResponseHandler[nethttp.Response] {
	return r.handler
}

func (r *nativeResponseHandlerOption) SetResponseHandler(responseHandler ResponseHandler[nethttp.Response]) {
	r.handler = responseHandler
}

func (r *nativeResponseHandlerOption) GetKey() ResponseOptionKey {
	return ResponseHandlerOptionKey
}
