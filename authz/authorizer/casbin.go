package authorizer

import (
	"github.com/casbin/casbin"
	"github.com/ideahitme/k8s-api-webhook/authz/unversioned"
)

/**
Authorizer enabled via casbin style file and using casbin package (github.com/casbin/casbin)
*/

// CasbinResource type authorizer for resources
type CasbinResource struct {
	casbin *casbin.Enforcer
}

// CasbinNonResource type authorizer for non-resources
type CasbinNonResource struct {
	casbin *casbin.Enforcer
}

const (
	casbinConfigPath = ""
	casbinModelPath  = ""
)

// NewCasbinResource reads csv policy file and constructs required policy enforcer
func NewCasbinResource(policyFile string) (*CasbinResource, error) {
	return nil, nil
}

// NewCasbinNonResource reads csv policy file and constructs required policy enforcer
func NewCasbinNonResource(policyFile string) (*CasbinNonResource, error) {
	return nil, nil
}

// IsAuthorized returns true, nil if the user is allowed to access specified
// resource object
func (c *CasbinResource) IsAuthorized(*unversioned.UserSpec, *unversioned.ResourceSpec) (bool, error) {
	return false, nil
}

// IsAuthorized returns true, nil if the user is allowed to access specified
// non resource object
func (c *CasbinNonResource) IsAuthorized(*unversioned.UserSpec, *unversioned.NonResourceSpec) (bool, error) {
	return false, nil
}
