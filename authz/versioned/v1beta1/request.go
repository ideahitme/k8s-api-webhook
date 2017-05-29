package v1beta1

import (
	"encoding/json"
	"io"

	"github.com/ideahitme/k8s-api-webhook/authz/unversioned"
)

/**
Incoming request to resource objects

{
  "apiVersion": "authorization.k8s.io/v1beta1",
  "kind": "SubjectAccessReview",
  "spec": {
    "resourceAttributes": {
      "namespace": "kittensandponies",
      "verb": "get",
      "group": "unicorn.example.org",
      "resource": "pods"
    },
    "user": "jane",
    "group": [
      "group1",
      "group2"
    ]
  }
}

request to non-resource objects:

{
  "apiVersion": "authorization.k8s.io/v1beta1",
  "kind": "SubjectAccessReview",
  "spec": {
    "nonResourceAttributes": {
      "path": "/debug",
      "verb": "get"
    },
    "user": "jane",
    "group": [
      "group1",
      "group2"
    ]
  }
}

/api, /apis, /metrics, /resetMetrics, /logs, /debug, /healthz, /swagger-ui/, /swaggerapi/, /ui, and /version. Clients require access to /api, /api/*, /apis, /apis/*, and /version
to discover what resources and versions are present on the server.
*/

// SubjectAccessReview checks whether or not a user or group can perform an action.
type SubjectAccessReview struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	// Spec holds information about the request being evaluated
	Spec SubjectAccessReviewSpec `json:"spec"`
}

// ResourceAttributes includes the authorization attributes available for resource requests to the Authorizer interface
type ResourceAttributes struct {
	// Namespace is the namespace of the action being requested.  Currently, there is no distinction between no namespace and all namespaces
	Namespace string `json:"namespace,omitempty"`
	// Verb is a kubernetes resource API verb, like: get, list, watch, create, update, delete, proxy.  "*" means all.
	Verb string `json:"verb,omitempty"`
	// Group is the API Group of the Resource.  "*" means all.
	Group string `json:"group,omitempty"`
	// Version is the API Version of the Resource.  "*" means all.
	Version string `json:"version,omitempty"`
	// Resource is one of the existing resource types.  "*" means all.
	Resource string `json:"resource,omitempty"`
	// Subresource is one of the existing resource types.  "" means none.
	Subresource string `json:"subresource,omitempty"`
	// Name is the name of the resource being requested for a "get" or deleted for a "delete". "" (empty) means all.
	Name string `json:"name,omitempty"`
}

// NonResourceAttributes includes the authorization attributes available for non-resource requests to the Authorizer interface
type NonResourceAttributes struct {
	// Path is the URL path of the request
	Path string `json:"path,omitempty"`
	// Verb is the standard HTTP verb
	Verb string `json:"verb,omitempty"`
}

// SubjectAccessReviewSpec is a description of the access request.  Exactly one of ResourceAuthorizationAttributes
// and NonResourceAuthorizationAttributes must be set
type SubjectAccessReviewSpec struct {
	// ResourceAuthorizationAttributes describes information for a resource access request
	ResourceAttributes *ResourceAttributes `json:"resourceAttributes,omitempty"`
	// NonResourceAttributes describes information for a non-resource access request
	NonResourceAttributes *NonResourceAttributes `json:"nonResourceAttributes,omitempty"`
	// User is the user you're testing for.
	User string `json:"user,omitempty"`
	// Groups is the groups you're testing for.
	Groups []string `json:"group,omitempty"`
	// Extra map[string][]string `json:"extra,omitempty"` ignore extra fields
}

// RequestParser implements extraction of the spec according to the official requirement for v1beta1 version
type RequestParser struct {
	body *SubjectAccessReview
}

// ReadBody parses the request body into k8s API Server specified format
// it should be called before extracting specs
func (req *RequestParser) ReadBody(body io.ReadCloser) error {
	decoder := json.NewDecoder(body)
	requestBody := &SubjectAccessReview{}

	err := decoder.Decode(requestBody)
	if err != nil {
		return err
	}
	defer body.Close()

	req.body = requestBody
	return nil
}

// IsResourceRequest returns true if the request is targeted for a resource
// for example "create pod"
func (req *RequestParser) IsResourceRequest() bool {
	return req.body.Spec.ResourceAttributes != nil
}

// IsNonResourceRequest returns true if the request is targeted for a non-resource,
// for example "metrics" exposed by API Server
func (req *RequestParser) IsNonResourceRequest() bool {
	return req.body.Spec.NonResourceAttributes != nil
}

// ExtractResourceSpecs extracts resource related fields from previously read request body
// see ReadBody method
func (req *RequestParser) ExtractResourceSpecs() *unversioned.ResourceSpec {
	if !req.IsResourceRequest() {
		return nil
	}
	return &unversioned.ResourceSpec{
		Namespace: req.body.Spec.ResourceAttributes.Namespace,
		Verb:      req.body.Spec.ResourceAttributes.Verb,
		Resource:  req.body.Spec.ResourceAttributes.Resource,
	}
}

// ExtractNonResourceSpecs extracts non-resource related fields from previously read request body
// see ReadBody method
func (req *RequestParser) ExtractNonResourceSpecs() *unversioned.NonResourceSpec {
	if !req.IsNonResourceRequest() {
		return nil
	}
	return &unversioned.NonResourceSpec{
		Path: req.body.Spec.NonResourceAttributes.Path,
		Verb: req.body.Spec.NonResourceAttributes.Verb,
	}
}

// ExtractUserSpecs reads the request body received from API server and extracts all required scopes by the user
func (req *RequestParser) ExtractUserSpecs() *unversioned.UserSpec {
	return &unversioned.UserSpec{
		User:   req.body.Spec.User,
		Groups: req.body.Spec.Groups,
	}
}
