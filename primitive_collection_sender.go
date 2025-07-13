package abstractions

import (
	"context"
	nethttp "net/http"
)

type PrimitiveCollectionSender interface {
	// SendPrimitiveCollection executes the HTTP request specified by the given nethttp.Request and returns the deserialized primitive response model collection.
	SendPrimitiveCollection(context context.Context, request nethttp.Request, typeName string, errorMappings ErrorMappings) ([]any, error)
}
