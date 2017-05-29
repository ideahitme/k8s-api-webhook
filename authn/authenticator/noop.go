package authenticator

import "github.com/ideahitme/k8s-api-webhook/authn/unversioned"

// Noop is the default authenticator which only returns empty user
type Noop struct{}

// Authenticate return empty userinfo
func (n Noop) Authenticate(token string) (*unversioned.UserInfo, error) {
	return nil, nil
}
