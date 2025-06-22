package authentication

import (
	"context"
)

// AuthenticationProvider authenticates the request.
type AuthenticationProvider[T any] interface {
	// AuthenticateRequest authenticates the provided request.
	AuthenticateRequest(context context.Context, request T, additionalAuthenticationContext map[string]interface{}) error
}
