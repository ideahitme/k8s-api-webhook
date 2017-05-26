package v1beta1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFailResponse(t *testing.T) {
	reason := "user does not have read access to the namespace"
	expectedRes := []byte(`{"apiVersion":"authorization.k8s.io/v1beta1","kind":"SubjectAccessReview","status":{"allowed":false,"reason":"` + reason + `"}}`)

	assert.Equal(t, expectedRes, ResponseConstructor{}.NewFailResponse(reason))
}

func TestNewSuccessResponse(t *testing.T) {
	expectedRes := []byte(`{"apiVersion":"authorization.k8s.io/v1beta1","kind":"SubjectAccessReview","status":{"allowed":true}}`)

	assert.Equal(t, expectedRes, ResponseConstructor{}.NewSuccessResponse())
}
