package main

import (
	"testing"

	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"
	httpHelper "github.com/easysoft/zentaoatf/cmd/test/helper/http"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestTokenApi(t *testing.T) {
	suite.RunSuite(t, new(TokenApiSuite))
}

type TokenApiSuite struct {
	suite.Suite
}

func (s *TokenApiSuite) BeforeEach(t provider.T) {
	commonTestHelper.ReplaceLabel(t, "TokenApi")
}

func (s *TokenApiSuite) TestTokenApi(t provider.T) {
	t.ID("7637")
	token := httpHelper.Login()

	t.Require().Greater(len(token), 6, "login failed")
}
