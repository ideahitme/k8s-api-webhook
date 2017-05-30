package authorizer

import (
	"io/ioutil"
	"os"

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

// reloader automatically reloads configuration files to reflect latest changes
// type reloader struct {
// 	policyFile string
// 	modelFile  string
// 	reloadable reloadable
// }

// type reloadable interface {
// 	LoadModel()
// 	LoadPolicy()
// }

// func (r reloader) Run() {
// 	for {
// 		r.reloadable.LoadModel()
// 		r.reloadable.LoadPolicy()
// 		time.Sleep(5 * time.Minute)
// 	}
// }

// NewCasbinResource reads csv policy file and constructs required policy enforcer
func NewCasbinResource(policyFile, modelFile string) (authz *CasbinResource, err error) {
	conf, err := GenerateCasbinConfigFile(policyFile, modelFile)
	if err != nil {
		return nil, err
	}
	defer os.Remove(conf.Name())
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			authz = nil
		}
	}()

	authz = &CasbinResource{
		casbin.NewEnforcer(conf.Name()),
	}

	return authz, nil
}

// NewCasbinNonResource reads csv policy file and constructs required policy enforcer
func NewCasbinNonResource(policyFile, modelFile string) (authz *CasbinNonResource, err error) {
	conf, err := GenerateCasbinConfigFile(policyFile, modelFile)
	if err != nil {
		return nil, err
	}
	defer os.Remove(conf.Name())
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			authz = nil
		}
	}()

	authz = &CasbinNonResource{
		casbin.NewEnforcer(conf.Name()),
	}

	return authz, nil
}

// IsAuthorized returns true, nil if the user is allowed to access specified
// resource object
func (c *CasbinResource) IsAuthorized(*unversioned.UserSpec, *unversioned.ResourceSpec) (allowed bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			allowed = false
		}
	}()
	allowed = c.casbin.Enforce()
	return
}

// IsAuthorized returns true, nil if the user is allowed to access specified
// non resource object
func (c *CasbinNonResource) IsAuthorized(*unversioned.UserSpec, *unversioned.NonResourceSpec) (allowed bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			allowed = false
		}
	}()
	allowed = c.casbin.Enforce()
	return
}

// GenerateCasbinConfigFile generates configuration file
func GenerateCasbinConfigFile(policyFile, modelFile string) (*os.File, error) {
	f, err := ioutil.TempFile("", "casbin-configuration-file")
	if err != nil {
		return nil, err
	}

	_, err = f.WriteString(`[default]
model_path = ` + modelFile + `

policy_backend = file

[file]
policy_path = ` + policyFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f, nil
}
