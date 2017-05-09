package v1beta1

import (
	"encoding/json"
	"fmt"
)

/**
Request body sample
{
	"apiVersion": "authentication.k8s.io/v1beta1",
	"kind": "TokenReview",
	"spec": {
		"token": "(BEARERTOKEN)"
	}
}
*/

// requestBody request proxied to webhook for authentication
type requestBodyType struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Spec       struct {
		Token string `json:"token"`
	} `json:"spec"`
}

// ParseToken retrieves token from the payload and validates request struct
func ParseToken(payload []byte) (string, error) {
	var body requestBodyType
	err := json.Unmarshal(payload, &body)
	if err != nil {
		return "", err
	}
	if body.APIVersion != "authentication.k8s.io/v1beta1" || body.Kind != "TokenReview" {
		return "", fmt.Errorf("invalid authn request format")
	}
	if body.Spec.Token == "" {
		return "", fmt.Errorf("token is missing")
	}
	return body.Spec.Token, nil
}
