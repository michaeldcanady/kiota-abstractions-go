package abstractions

import (
	"context"
	nethttp "net/http"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EnumSender interface {
	// SendEnum executes the HTTP request specified by the given nethttp.Request and returns the deserialized response model.
	SendEnum(context context.Context, request nethttp.Request, parser s.EnumFactory, errorMappings ErrorMappings) (any, error)
}
