package authorizer

import (
	"testing"

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

	suite.casbinAuthz, err = NewCasbinResource("")
	suite.NoError(err, "file should exist!")
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

	suite.casbinAuthz, err = NewCasbinNonResource("")
	suite.NoError(err, "file should exist!")
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
