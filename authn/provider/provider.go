package provider

import "github.com/ideahitme/k8s-api-webhook/authn/unversioned"

// Authenticator is the general concept providing authentication information
// Authenticate method should return error only if internal error happened
// if user could not be identified returned value should be nil, nil
type Authenticator interface {
	Authenticate(token string) (*unversioned.UserInfo, error)
}
