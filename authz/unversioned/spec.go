package unversioned

// UserSpec defines information known regarding the user
type UserSpec struct {
	User   string
	Groups []string
}

// ResourceSpec defines the specifications of the authorization request
// that would include resource and the user making the request
type ResourceSpec struct {
	Namespace string
	Verb      string
	Resource  string
}

// NonResourceSpec defines the specifications of the authorization request
// that would non-resource and the user making the request
type NonResourceSpec struct {
	Path string
	Verb string
}
