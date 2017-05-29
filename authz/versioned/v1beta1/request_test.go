package v1beta1

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RequestParserTestSuite struct {
	suite.Suite
	validPayload       string
	invalidPayload     string
	resourceRequest    *RequestParser
	nonResourceRequest *RequestParser
}

func (suite *RequestParserTestSuite) SetupTest() {
	suite.validPayload = `{
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
	suite.invalidPayload = `
	{
		"apiVersion": "v1",
	}
	`
	resourceParser := RequestParser{}
	body := ioutil.NopCloser(strings.NewReader(suite.validPayload))
	err := resourceParser.ReadBody(body)
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.resourceRequest = &resourceParser

	nonResourcePayload := `
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

	nonResourceParser := RequestParser{}
	body = ioutil.NopCloser(strings.NewReader(nonResourcePayload))
	err = nonResourceParser.ReadBody(body)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.nonResourceRequest = &nonResourceParser
}

func (suite *RequestParserTestSuite) TestReadValidBody() {
	parser := RequestParser{}
	body := ioutil.NopCloser(strings.NewReader(suite.validPayload))
	err := parser.ReadBody(body)
	suite.NoError(err, "should not return error")
}

func (suite *RequestParserTestSuite) TestReadInvalidBody() {
	parser := RequestParser{}
	body := ioutil.NopCloser(strings.NewReader(suite.invalidPayload))
	err := parser.ReadBody(body)
	suite.Error(err, "should return error")
}

func (suite *RequestParserTestSuite) TestIsResourceRequest() {
	suite.True(suite.resourceRequest.IsResourceRequest(), "should return true")
	suite.False(suite.nonResourceRequest.IsResourceRequest(), "should return false")
}

func (suite *RequestParserTestSuite) TestIsNonResourceRequest() {
	suite.False(suite.resourceRequest.IsNonResourceRequest(), "should return true")
	suite.True(suite.nonResourceRequest.IsNonResourceRequest(), "should return false")
}

func (suite *RequestParserTestSuite) TestExtractResourceSpecs() {
	suite.Nil(suite.nonResourceRequest.ExtractResourceSpecs(), "should return nil")
	suite.NotNil(suite.resourceRequest.ExtractResourceSpecs(), "should not return nil")
}

func (suite *RequestParserTestSuite) TestExtractNonResourceSpecs() {
	suite.Nil(suite.resourceRequest.ExtractNonResourceSpecs(), "should return nil")
	suite.NotNil(suite.nonResourceRequest.ExtractNonResourceSpecs(), "should not return nil")
}

func (suite *RequestParserTestSuite) TestExtractUserSpecs() {
	suite.Equal("jane", suite.resourceRequest.ExtractUserSpecs().User, "should be equal")
	suite.Equal([]string{"group1", "group2"}, suite.resourceRequest.ExtractUserSpecs().Groups, "should be equal")
}

func TestRequestParser(t *testing.T) {
	suite.Run(t, new(RequestParserTestSuite))
}
