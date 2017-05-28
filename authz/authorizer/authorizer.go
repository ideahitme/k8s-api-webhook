package authorizer

import (
	"github.com/ideahitme/k8s-api-webhook/authz/unversioned"
)

// ResourceAuthorizer interface to be implemented by authorizer based on resource
type ResourceAuthorizer interface {
	IsAuthorized(*unversioned.UserSpec, *unversioned.ResourceSpec) (bool, error)
}

// NonResourceAuthorizer interface to be implemented by authorizer based on non-resources
type NonResourceAuthorizer interface {
	IsAuthorized(*unversioned.UserSpec, *unversioned.NonResourceSpec) (bool, error)
}
