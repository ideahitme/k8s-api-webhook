package authorizer

import "github.com/ideahitme/k8s-api-webhook/authz/unversioned"

// ResourceUnauthorizer is authorizer which does not allow access to any of the resources
type ResourceUnauthorizer struct{}

// NonResourceUnauthorizer is authorizer which does not allow access to any of the non-resources
type NonResourceUnauthorizer struct{}

// IsAuthorized always returns false
func (ResourceUnauthorizer) IsAuthorized(*unversioned.UserSpec, *unversioned.ResourceSpec) (bool, error) {
	return false, nil
}

// IsAuthorized always returns false
func (NonResourceUnauthorizer) IsAuthorized(*unversioned.UserSpec, *unversioned.NonResourceSpec) (bool, error) {
	return false, nil
}
