package v1beta1

import (
	"encoding/json"
	"fmt"
	"io"
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

// RequestParser implements parsing of the request for v1beta1 version
type RequestParser struct {
}

// ExtractToken retrieves token from the payload and validates request struct
func (r RequestParser) ExtractToken(reqBody io.ReadCloser) (string, error) {
	decoder := json.NewDecoder(reqBody)
	var body requestBodyType
	err := decoder.Decode(&body)
	if err != nil {
		return "", err
	}
	defer reqBody.Close()

	if body.APIVersion != "authentication.k8s.io/v1beta1" || body.Kind != "TokenReview" {
		return "", fmt.Errorf("invalid authn request format")
	}
	if body.Spec.Token == "" {
		return "", fmt.Errorf("token is missing")
	}
	return body.Spec.Token, nil
}
