package v1beta1

import (
	"testing"

	"github.com/ideahitme/k8s-api-webhook/authn/unversioned"
	"github.com/stretchr/testify/assert"
)

func TestNewFailResponse(t *testing.T) {
	constr := ResponseConstructor{}
	assert.Equal(t, constr.NewFailResponse(), []byte(`
	{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kind": "TokenReview",
		"status": {
			"authenticated": false
		}
	}
`))
}

func TestNewSuccessResponse(t *testing.T) {
	for _, ti := range []struct {
		title            string
		user             *unversioned.UserInfo
		expectedResponse []byte
	}{
		{
			title:            "empty user",
			user:             nil,
			expectedResponse: nil,
		},
		{
			title: "non-empty user",
			user: &unversioned.UserInfo{
				UID:    "user-uid",
				Name:   "foo",
				Groups: []string{"bar", "baz"},
			},
			expectedResponse: []byte(`{"apiVersion":"authentication.k8s.io/v1beta1","kind":"TokenReview","status":{"authenticated":true,"user":{"username":"foo","uid":"user-uid","groups":["bar","baz"],"extra":{}}}}`),
		},
	} {
		t.Run(ti.title, func(t *testing.T) {
			constr := ResponseConstructor{}
			res := constr.NewSuccessResponse(ti.user)
			assert.Equal(t, ti.expectedResponse, res)
		})
	}
}
