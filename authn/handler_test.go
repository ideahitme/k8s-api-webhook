package authn

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ideahitme/k8s-api-webhook/authn/provider"
	"github.com/ideahitme/k8s-api-webhook/authn/unversioned"
	"github.com/ideahitme/k8s-api-webhook/authn/v1beta1"
	"github.com/ideahitme/k8s-api-webhook/internal/testutils"
)

func TestNewAuthenticationHandler(t *testing.T) {
	h := NewAuthenticationHandler(provider.StaticAuthenticator{}, WithAPIVersion(V1Beta1))
	assert.Equal(t, h.authProvider, provider.StaticAuthenticator{})
}

func TestWithAPIVersion(t *testing.T) {
	h := &AuthenticationHandler{}
	WithAPIVersion(V1Beta1)(h)
	assert.Equal(t, h.reqParser, v1beta1.RequestParser{})
	assert.Equal(t, h.resConstructor, v1beta1.ResponseConstructor{})

	WithAPIVersion(-1)(h)
	assert.Equal(t, h.reqParser, v1beta1.RequestParser{})
	assert.Equal(t, h.resConstructor, v1beta1.ResponseConstructor{})
}

func TestServeHTTP(t *testing.T) {
	f := testutils.GenerateTestData([][]string{
		{
			"token-1", "foo", "42",
		},
		{
			"token-2", "bar", "99", "Admin", "Owner",
		},
	})
	defer os.Remove(f.Name())

	staticAuthenticator, err := provider.NewStaticAuthenticator(f.Name())
	assert.Nil(t, err)
	failAuthenticator := errAuthenticator{invalidToken: "cause-error"}

	handler := NewAuthenticationHandler(provider.NewAggregator(staticAuthenticator, failAuthenticator))
	mockServer := httptest.NewServer(handler)
	defer mockServer.Close()

	url := mockServer.URL

	var (
		res  *http.Response
		body []byte
	)

	// token cannot be identified - unauthorized request
	res, err = sendRequest(url, []byte(`{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kind": "TokenReview",
		"spec": {
			"token": "my-token"
		}
	}`))

	body, err = ioutil.ReadAll(res.Body)

	assert.Nil(t, err)
	assert.Equal(t, v1beta1.ResponseConstructor{}.NewFailResponse(), body)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// not json payload
	res, err = sendRequest(url, []byte(`{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kin`))

	body, err = ioutil.ReadAll(res.Body)

	assert.Nil(t, err)
	assert.Equal(t, v1beta1.ResponseConstructor{}.NewFailResponse(), body)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	// valid request
	res, err = sendRequest(url, []byte(`{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kind": "TokenReview",
		"spec": {
			"token": "token-2"
		}
	}`))

	body, err = ioutil.ReadAll(res.Body)

	assert.Nil(t, err)
	assert.Equal(t, v1beta1.ResponseConstructor{}.NewSuccessResponse(
		&unversioned.UserInfo{
			UID:    "99",
			Name:   "bar",
			Groups: []string{"Admin", "Owner"},
		},
	), body)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// error in authenticator
	res, err = sendRequest(url, []byte(`{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kind": "TokenReview",
		"spec": {
			"token": "cause-error"
		}
	}`))

	body, err = ioutil.ReadAll(res.Body)

	assert.Nil(t, err)
	assert.Equal(t, v1beta1.ResponseConstructor{}.NewFailResponse(), body)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func sendRequest(url string, payload []byte) (*http.Response, error) {
	return http.Post(url, "application/json", bytes.NewReader(payload))
}

type errAuthenticator struct {
	invalidToken string
}

func (ma errAuthenticator) Authenticate(token string) (*unversioned.UserInfo, error) {
	if token == ma.invalidToken {
		return nil, fmt.Errorf("failed to authenticate")
	}
	return nil, nil
}
