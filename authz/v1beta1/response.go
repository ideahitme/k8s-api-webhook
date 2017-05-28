package v1beta1

import (
	"encoding/json"
)

/**

Successfuly response:
{
  "apiVersion": "authorization.k8s.io/v1beta1",
  "kind": "SubjectAccessReview",
  "status": {
    "allowed": true
  }
}

Unauthorized response:
{
  "apiVersion": "authorization.k8s.io/v1beta1",
  "kind": "SubjectAccessReview",
  "status": {
    "allowed": false,
    "reason": "user does not have read access to the namespace"
  }
}

*/

type response struct {
	APIVersion string         `json:"apiVersion"`
	Kind       string         `json:"kind"`
	Status     responseStatus `json:"status"`
}

type responseStatus struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

// ResponseConstructor constructs response to authorization requests according to official requirements of v1beta1 version
type ResponseConstructor struct {
}

// NewSuccessResponse returns response to API server to signify that user is authorized
func (ResponseConstructor) NewSuccessResponse() []byte {
	res := response{
		APIVersion: "authorization.k8s.io/v1beta1",
		Kind:       "SubjectAccessReview",
		Status: responseStatus{
			Allowed: true,
		},
	}
	b, _ := json.Marshal(res)
	return b
}

// NewFailResponse returns response to API server to signify that user is not authorized
func (ResponseConstructor) NewFailResponse(reason string) []byte {
	res := response{
		APIVersion: "authorization.k8s.io/v1beta1",
		Kind:       "SubjectAccessReview",
		Status: responseStatus{
			Allowed: false,
			Reason:  reason,
		},
	}
	b, _ := json.Marshal(res)
	return b
}
