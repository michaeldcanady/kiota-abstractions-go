package abstractions

import (
	"context"
	nethttp "net/http"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CollectionSender interface {
	// SendCollection executes the HTTP request specified by the given nethttp.Request and returns the deserialized response model collection.
	SendCollection(context context.Context, request nethttp.Request, constructor s.ParsableFactory, errorMappings ErrorMappings) ([]s.Parsable, error)
}
