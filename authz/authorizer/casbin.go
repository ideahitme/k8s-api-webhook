package authorizer

import (
	"github.com/casbin/casbin"
)

/**
Authorizer enabled via casbin style file and using casbin package (github.com/casbin/casbin)
*/

// Casbin type authorizer
type Casbin struct {
	casbin *casbin.Enforcer
}

const (
	casbinConfigPath = ""
	casbinModelPath  = ""
)

// NewCasbin reads csv policy file and constructs required policy enforcer
func NewCasbin(policyFile string) (*Casbin, error) {

	return nil, nil
}

// ResourceEnforce returns true, nil if the user is allowed to access specified
// resource object
func (c *Casbin) ResourceEnforce() (bool, error) {
	return false, nil
}

// NonResourceEnforce returns true, nil if the user is allowed to access specified
// non resource object
func (c *Casbin) NonResourceEnforce() (bool, error) {
	return false, nil
}
