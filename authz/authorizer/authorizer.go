package authorizer

// Authorizer defines the interface required to handle access to resource/non-resource endpoints of k8s API server
type Authorizer interface {
	ResourceEnforce() (allowed bool, err error)
	NonResourceEnforce() (allowed bool, err error)
}
