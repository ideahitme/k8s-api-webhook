package v1beta1

import (
	"encoding/json"

	"github.com/ideahitme/k8s-api-webhook/authn/unversioned"
)

/**

Successful authentication response:
{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "status": {
    "authenticated": true,
    "user": {
      "username": "janedoe@example.com",
      "uid": "42",
      "groups": [
        "developers",
        "qa"
      ],
      "extra": {
        "extrafield1": [
          "extravalue1",
          "extravalue2"
        ]
      }
    }
  }
}

Failed authentication response:
{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "status": {
    "authenticated": false
  }
}

*/

type successResponseStatus struct {
	Authenticated bool                `json:"authenticated"`
	User          successResponseUser `json:"user"`
}

type successResponseUser struct {
	Username string            `json:"username"`
	UID      string            `json:"uid"`
	Groups   []string          `json:"groups"`
	Extra    map[string]string `json:"extra"`
}

type successResponseBody struct {
	APIVersion string                `json:"apiVersion"`
	Kind       string                `json:"kind"`
	Status     successResponseStatus `json:"status"`
}

// NewSuccessResponse returns the typical success response for v1beta1 APIVersion
func NewSuccessResponse(user *unversioned.UserInfo) []byte {
	resp := &successResponseBody{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind:       "TokenReview",
		Status: successResponseStatus{
			Authenticated: true,
			User: successResponseUser{
				Username: user.Name,
				UID:      user.UID,
				Groups:   user.Groups,
				Extra:    map[string]string{},
			},
		},
	}
	b, _ := json.Marshal(resp)
	return b
}

// NewFailResponse returns the typical fail response for v1beta1 APIVersion
func NewFailResponse() []byte {
	return []byte(`
	{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kind": "TokenReview",
		"status": {
			"authenticated": false
		}
	}
`)
}
