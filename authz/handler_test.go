package authz

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/ideahitme/k8s-api-webhook/authz/authorizer"
	"github.com/ideahitme/k8s-api-webhook/authz/unversioned"
	"github.com/ideahitme/k8s-api-webhook/authz/versioned"
	"github.com/ideahitme/k8s-api-webhook/authz/versioned/v1beta1"
)

type AuthorizationHandlerSuite struct {
	suite.Suite
	unauthorizeEndpoint      string
	mockEndpoint             string
	hacker                   *unversioned.UserSpec
	legitUser                *unversioned.UserSpec
	resourcePayload          string
	nonResourcePayload       string
	hackerResourcePayload    string
	hackerNonResourcePayload string
}

func (suite *AuthorizationHandlerSuite) SetupTest() {
	unauthzServer := httptest.NewServer(CreateAuthorizationHandler())
	suite.unauthorizeEndpoint = unauthzServer.URL

	suite.resourcePayload = `{
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
	}`
	suite.nonResourcePayload = `
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
	`
	suite.hackerResourcePayload = `{
		"apiVersion": "authorization.k8s.io/v1beta1",
		"kind": "SubjectAccessReview",
		"spec": {
			"resourceAttributes": {
				"namespace": "kittensandponies",
				"verb": "get",
				"group": "unicorn.example.org",
				"resource": "pods"
			},
			"user": "hacker"
		}
	}`
	suite.hackerNonResourcePayload = `
	{
		"apiVersion": "authorization.k8s.io/v1beta1",
		"kind": "SubjectAccessReview",
		"spec": {
			"nonResourceAttributes": {
				"path": "/debug",
				"verb": "get"
			},
			"user": "hacker"
		}
	}
	`

	suite.hacker = &unversioned.UserSpec{
		User: "hacker",
	}
	suite.legitUser = &unversioned.UserSpec{
		User:   "jane",
		Groups: []string{"group1", "group2"},
	}

	//configure mock resource handlers
	mockResourceHandler := new(MockResourceAuthorizer)
	mockResourceHandler.On("IsAuthorized", suite.legitUser, &unversioned.ResourceSpec{
		Namespace: "kittensandponies",
		Verb:      "get",
		Resource:  "pods",
	}).Return(true, nil)
	mockResourceHandler.On("IsAuthorized", suite.hacker, &unversioned.ResourceSpec{
		Namespace: "kittensandponies",
		Verb:      "get",
		Resource:  "pods",
	}).Return(false, errors.New("hackers not allowed"))

	mockNonResourceHandler := new(MockNonResourceAuthorizer)
	mockNonResourceHandler.On("IsAuthorized", suite.legitUser, &unversioned.NonResourceSpec{
		Path: "/debug",
		Verb: "get",
	}).Return(true, nil)
	mockNonResourceHandler.On("IsAuthorized", suite.hacker, &unversioned.NonResourceSpec{
		Path: "/debug",
		Verb: "get",
	}).Return(false, errors.New("hackers not allowed"))

	mockHandler := CreateAuthorizationHandler().
		WithNonResourceAuthorizer(mockNonResourceHandler).
		WithResourceAuthorizer(mockResourceHandler)
	mockServer := httptest.NewServer(mockHandler)
	suite.mockEndpoint = mockServer.URL
}

func (suite *AuthorizationHandlerSuite) TestCreateAuthorizationHandler() {
	h := CreateAuthorizationHandler()
	suite.IsType(authorizer.ResourceUnauthorizer{}, h.resourceAuthorizer, "default should be unauthorizer")
	suite.IsType(authorizer.NonResourceUnauthorizer{}, h.nonResourceAuthorizer, "default should be unauthorizer")
	suite.IsType(&v1beta1.ResponseConstructor{}, h.resConstructor, "default should be v1beta1")
	suite.IsType(&v1beta1.RequestParser{}, h.reqParser, "default should be v1beta1")
}

// TestExtensions makes sure with chaining works as expected
func (suite *AuthorizationHandlerSuite) TestExtensions() {
	h := CreateAuthorizationHandler().WithAPIVersion(versioned.V1Beta1).WithNonResourceAuthorizer(&authorizer.CasbinNonResource{}).
		WithResourceAuthorizer(&authorizer.CasbinResource{})
	suite.IsType(&authorizer.CasbinResource{}, h.resourceAuthorizer, "default should be overridden with casbin")
	suite.IsType(&authorizer.CasbinNonResource{}, h.nonResourceAuthorizer, "default should be overridden with casbin")
	suite.IsType(&v1beta1.ResponseConstructor{}, h.resConstructor, "default should be v1beta1")
	suite.IsType(&v1beta1.RequestParser{}, h.reqParser, "default should be v1beta1")
}

// TestUnauthorized is a test againt a default authz endpoints
func (suite *AuthorizationHandlerSuite) TestUnauthorized() {
	res, err := SendRequest(suite.unauthorizeEndpoint, []byte(`rubbish`))
	suite.NoError(err, "should not return error, server is up")
	suite.Equal(http.StatusBadRequest, res.StatusCode, "should return bad request")

	//this time no malformed request
	res, err = SendRequest(suite.unauthorizeEndpoint, []byte(suite.resourcePayload))
	suite.NoError(err, "should not return error, server is up")
	suite.Equal(http.StatusUnauthorized, res.StatusCode, "should return unauthorized")

	res, err = SendRequest(suite.unauthorizeEndpoint, []byte(suite.nonResourcePayload))
	suite.NoError(err, "should not return error, server is up")
	suite.Equal(http.StatusUnauthorized, res.StatusCode, "should return unauthorized")
}

// TestWithMockAuthorizer is a test againt a default authz endpoints
func (suite *AuthorizationHandlerSuite) TestWithMockAuthorizer() {
	res, err := SendRequest(suite.mockEndpoint, []byte(suite.resourcePayload))
	suite.NoError(err, "should not return error, server is up")
	suite.Equal(http.StatusOK, res.StatusCode, "should return user authorized")

	res, err = SendRequest(suite.mockEndpoint, []byte(suite.nonResourcePayload))
	suite.NoError(err, "should not return error, server is up")
	suite.Equal(http.StatusOK, res.StatusCode, "should return user authorized")

	//passing hacker payload should cause an error
	res, err = SendRequest(suite.mockEndpoint, []byte(suite.hackerResourcePayload))
	suite.NoError(err, "should not return error, server is up")
	suite.Equal(http.StatusInternalServerError, res.StatusCode, "should return internal server error")

	res, err = SendRequest(suite.mockEndpoint, []byte(suite.hackerNonResourcePayload))
	suite.NoError(err, "should not return error, server is up")
	suite.Equal(http.StatusInternalServerError, res.StatusCode, "should return internal server error")
}

func TestAuthorizationHandler(t *testing.T) {
	suite.Run(t, new(AuthorizationHandlerSuite))
}

func SendRequest(endpoint string, payload []byte) (*http.Response, error) {
	return http.Post(endpoint, "application/json", bytes.NewReader(payload))
}

type MockResourceAuthorizer struct {
	mock.Mock
}

func (m *MockResourceAuthorizer) IsAuthorized(user *unversioned.UserSpec, resourceSpec *unversioned.ResourceSpec) (bool, error) {
	args := m.Called(user, resourceSpec)
	return args.Bool(0), args.Error(1)
}

type MockNonResourceAuthorizer struct {
	mock.Mock
}

func (m *MockNonResourceAuthorizer) IsAuthorized(user *unversioned.UserSpec, nonResourceSpec *unversioned.NonResourceSpec) (bool, error) {
	args := m.Called(user, nonResourceSpec)
	return args.Bool(0), args.Error(1)
}
