package abstractions

import (
	"context"
	nethttp "net/http"
)

type PrimitiveSender interface {
	// SendPrimitive executes the HTTP request specified by the given nethttp.Request and returns the deserialized primitive response model.
	SendPrimitive(context context.Context, request nethttp.Request, typeName string, errorMappings ErrorMappings) (any, error)
}
