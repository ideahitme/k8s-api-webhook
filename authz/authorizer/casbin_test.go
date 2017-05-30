package authorizer

import (
	"testing"

	"io/ioutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var _ ResourceAuthorizer = &CasbinResource{}
var _ NonResourceAuthorizer = &CasbinNonResource{}

type CasbinResourceSuite struct {
	suite.Suite
	casbinAuthz *CasbinResource
}

func (suite *CasbinResourceSuite) SetupTest() {
	var err error

	suite.casbinAuthz, err = NewCasbinResource("testdata/resource/policy.csv", "testdata/resource/model.conf")
	suite.NoError(err, "casbin client could not be created!")
}

func (suite *CasbinResourceSuite) TestInitErrors() {
	authz, err := NewCasbinResource("", "")
	suite.Nil(authz, "should be nil")
	suite.Error(err, "error should be returned")
}

func (suite *CasbinResourceSuite) TestIsAuthorized() {
	allowed, err := suite.casbinAuthz.IsAuthorized(nil, nil)
	suite.False(allowed, "should not be allowed")
	suite.NoError(err, "no error should be returned")
}

type CasbinNonResourceSuite struct {
	suite.Suite
	casbinAuthz *CasbinNonResource
}

func (suite *CasbinNonResourceSuite) SetupTest() {
	var err error

	suite.casbinAuthz, err = NewCasbinNonResource("testdata/non-resource/policy.csv", "testdata/non-resource/model.conf")
	suite.NoError(err, "casbin client could not be created!")
}

func (suite *CasbinNonResourceSuite) TestInitErrors() {
	authz, err := NewCasbinNonResource("", "")
	suite.Nil(authz, "should be nil")
	suite.Error(err, "error should be returned")
}

func (suite *CasbinNonResourceSuite) TestIsAuthorized() {
	allowed, err := suite.casbinAuthz.IsAuthorized(nil, nil)
	suite.False(allowed, "should not be allowed")
	suite.NoError(err, "no error should be returned")
}

func TestCasbin(t *testing.T) {
	suite.Run(t, new(CasbinResourceSuite))
	suite.Run(t, new(CasbinNonResourceSuite))
}

func TestGenerateCasbinConfigFile(t *testing.T) {
	f, err := GenerateCasbinConfigFile("policy-file.csv", "model-file.conf")
	assert.NoError(t, err, "no error expected")
	defer f.Close()

	b, err := ioutil.ReadFile(f.Name())
	assert.NoError(t, err, "no error expected")

	assert.Equal(t, `[default]
model_path = model-file.conf

policy_backend = file

[file]
policy_path = policy-file.csv`, string(b))
}
