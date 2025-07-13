// Package abstractions provides the base infrastructure for the Kiota-generated SDKs to function.
// It defines multiple concepts related to abstract HTTP requests, serialization, and authentication.
// These concepts can then be implemented independently without tying the SDKs to any specific implementation.
// Kiota also provides default implementations for these concepts.
// Checkout:
// - github.com/microsoft/kiota/authentication/go/azure
// - github.com/microsoft/kiota/http/go/nethttp
// - github.com/microsoft/kiota/serialization/go/json
package abstractions

import (
	"context"

	"github.com/microsoft/kiota-abstractions-go/store"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ErrorMappings is a mapping of status codes to error types factories.
type ErrorMappings map[string]s.ParsableFactory

type Sender interface {
}

// RequestAdapter is the service responsible for translating abstract nethttp.Request into native HTTP requests.
type RequestAdapter interface {
	// GetSerializationWriterFactory returns the serialization writer factory currently in use for the request adapter service.
	GetSerializationWriterFactory() s.SerializationWriterFactory
	// EnableBackingStore enables the backing store proxies for the SerializationWriters and ParseNodes in use.
	EnableBackingStore(factory store.BackingStoreFactory)
	// ConvertToNativeRequest converts the given RequestInformation into a native HTTP request.
	ConvertToNativeRequest(context context.Context, requestInfo *RequestInformation) (any, error)
}
