package authentication

import (
	"context"
	"errors"
	"fmt"

	nethttp "net/http"
)

const (
	authorizationHeader = "Authorization"
	claimsKey           = "claims"
)

var _ AuthenticationProvider[*nethttp.Request] = (*BaseBearerTokenAuthenticationProvider)(nil)

// BaseBearerTokenAuthenticationProvider provides a base class implementing AuthenticationProvider for Bearer token scheme.
type BaseBearerTokenAuthenticationProvider struct {
	// accessTokenProvider is called by the BaseBearerTokenAuthenticationProvider class to authenticate the request via the returned access token.
	accessTokenProvider AccessTokenProvider
}

// NewBaseBearerTokenAuthenticationProvider creates a new instance of the BaseBearerTokenAuthenticationProvider class.
func NewBaseBearerTokenAuthenticationProvider(accessTokenProvider AccessTokenProvider) *BaseBearerTokenAuthenticationProvider {
	return &BaseBearerTokenAuthenticationProvider{accessTokenProvider}
}

// AuthenticateRequest authenticates the provided RequestInformation instance using the provided authorization token callback.
func (provider *BaseBearerTokenAuthenticationProvider) AuthenticateRequest(ctx context.Context, request *nethttp.Request, additionalAuthenticationContext map[string]interface{}) error {
	if request == nil {
		return errors.New("request is nil")
	}
	if request.Header == nil {
		request.Header = make(nethttp.Header)
	}
	if provider.accessTokenProvider == nil {
		return errors.New("this class needs to be initialized with an access token provider")
	}
	if _, ok := request.Header[authorizationHeader]; additionalAuthenticationContext[claimsKey] != nil && ok {
		request.Header.Del(authorizationHeader)
	}
	if _, ok := request.Header[authorizationHeader]; !ok {
		uri := request.URL
		token, err := provider.accessTokenProvider.GetAuthorizationToken(ctx, uri, additionalAuthenticationContext)
		if err != nil {
			return err
		}
		if token != "" {
			request.Header.Add(authorizationHeader, fmt.Sprintf("Bearer %s", token))
		}
	}

	return nil
}

// GetAuthorizationTokenProvider returns the access token provider the BaseBearerTokenAuthenticationProvider class uses to authenticate the request.
func (provider *BaseBearerTokenAuthenticationProvider) GetAuthorizationTokenProvider() AccessTokenProvider {
	return provider.accessTokenProvider
}
