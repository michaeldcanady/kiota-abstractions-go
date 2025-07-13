package abstractions

import (
	"context"
	nethttp "net/http"
)

type NoContentSender interface {
	// SendNoContent executes the HTTP request specified by the given nethttp.Request with no return content.
	SendNoContent(context context.Context, request nethttp.Request, errorMappings ErrorMappings) error
}
