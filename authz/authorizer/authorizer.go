package authorizer

import (
	"github.com/ideahitme/k8s-api-webhook/authz/unversioned"
)

// Authorizer defines the interface required to handle access to resource/non-resource endpoints of k8s API server
type Authorizer interface {
	ResourceEnforce(*unversioned.UserSpec, *unversioned.ResourceSpec) (allowed bool, err error)
	NonResourceEnforce(*unversioned.UserSpec, *unversioned.NonResourceSpec) (allowed bool, err error)
}
