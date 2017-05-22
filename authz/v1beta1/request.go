package v1beta1

import (
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

// RequestParser implements extraction of the spec according to the official requirement for v1beta1 version
type RequestParser struct {
}

// ExtractScope reads the request body received from API server and extracts all required scopes by the user
func (req RequestParser) ExtractScope(io.ReadCloser) (*unversioned.Scope, error) {
	return nil, nil
}
