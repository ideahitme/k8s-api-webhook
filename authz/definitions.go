package authz

import (
	"io"

	"github.com/ideahitme/k8s-api-webhook/authz/unversioned"
)

// APIVersion allows to select the version API server for integration
type APIVersion int

const (
	//V1Beta1 version of the API Server
	V1Beta1 APIVersion = iota
)

// ResponseConstructor provides an interface to apiVersion dependent request formats
type ResponseConstructor interface {
	NewFailResponse(reason string) []byte
	NewSuccessResponse() []byte
}

// RequestParser provides an interface to parse and retrieve the token from the authentication request
type RequestParser interface {
	ExtractScope(io.ReadCloser) (*unversioned.Scope, error)
}
