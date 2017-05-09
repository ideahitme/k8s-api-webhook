package provider

import (
	"github.com/ideahitme/k8s-api-webhook/authn/unversioned"
)

// Aggregator aggregates authenticators
type Aggregator struct {
	providers []Authenticator
}

// NewAggregator receives list of authenticators
// order of passed authenticator determines in which order they will be queried
func NewAggregator(authenticators ...Authenticator) Aggregator {
	return Aggregator{
		providers: authenticators,
	}
}

// Authenticate try all providers until first error or success hit
func (a Aggregator) Authenticate(token string) (*unversioned.UserInfo, error) {
	for _, authn := range a.providers {
		user, err := authn.Authenticate(token)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return user, nil
		}
	}
	return nil, nil
}
