package authentication

import (
	"context"
)

var _ AuthenticationProvider[any] = (*AnonymousAuthenticationProvider[any])(nil)

// AnonymousAuthenticationProvider implements the AuthenticationProvider interface does not perform any authentication.
type AnonymousAuthenticationProvider[T any] struct {
}

// AuthenticateRequest is a placeholder method that "authenticates" the RequestInformation instance: no-op.
func (provider *AnonymousAuthenticationProvider[T]) AuthenticateRequest(context context.Context, request T, additionalAuthenticationContext map[string]interface{}) error {
	return nil
}
