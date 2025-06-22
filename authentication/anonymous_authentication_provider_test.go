package authentication

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestAnonymousProviderHonoursInterface(t *testing.T) {
	instance := &AnonymousAuthenticationProvider[any]{}
	assert.Implements(t, (*AuthenticationProvider[any])(nil), instance)
}
