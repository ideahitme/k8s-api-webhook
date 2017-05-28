package authorizer

import (
	"github.com/casbin/casbin"
	"github.com/ideahitme/k8s-api-webhook/authz/unversioned"
)

/**
Authorizer enabled via casbin style file and using casbin package (github.com/casbin/casbin)
*/

// ResourceCasbin type authorizer for resources
type ResourceCasbin struct {
	casbin *casbin.Enforcer
}

// NonResourceCasbin type authorizer for non-resources
type NonResourceCasbin struct {
	casbin *casbin.Enforcer
}

const (
	casbinConfigPath = ""
	casbinModelPath  = ""
)

// NewCasbinResource reads csv policy file and constructs required policy enforcer
func NewCasbinResource(policyFile string) (*ResourceCasbin, error) {
	return nil, nil
}

// NewCasbinNonResource reads csv policy file and constructs required policy enforcer
func NewCasbinNonResource(policyFile string) (*NonResourceCasbin, error) {
	return nil, nil
}

// IsAuthorized returns true, nil if the user is allowed to access specified
// resource object
func (c *ResourceCasbin) IsAuthorized(*unversioned.UserSpec, *unversioned.ResourceSpec) (bool, error) {
	return false, nil
}

// IsAuthorized returns true, nil if the user is allowed to access specified
// non resource object
func (c *NonResourceCasbin) IsAuthorized(*unversioned.UserSpec, *unversioned.NonResourceSpec) (bool, error) {
	return false, nil
}
