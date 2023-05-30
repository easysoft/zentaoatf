package main

import (
	"github.com/bitly/go-simplejson"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestTokenApi(t *testing.T) {
	suite.RunSuite(t, new(TokenApiSuite))
}

type TokenApiSuite struct {
	suite.Suite
}

func (s *TokenApiSuite) BeforeEach(t provider.T) {
	t.ID("1")
	t.AddSubSuite("TokenApi")
}

func (s *TokenApiSuite) TestTokenApi(t provider.T) {
	login(t)

	t.Require().Greater(len(ZentaoToken), 6, "login failed")
}

func login(t provider.T) {
	url := zentaoHelper.GenApiUrl("tokens", nil, constTestHelper.ZtfUrl)

	params := map[string]string{
		"account":  constTestHelper.ZentaoUsername,
		"password": constTestHelper.ZentaoPassword,
	}
	bodyBytes, err := httpUtils.Post(url, params)
	if err != nil {
		return
	}

	jsn, err := simplejson.NewJson(bodyBytes)
	if err != nil || jsn == nil {
		return
	}

	mp, err := jsn.Map()
	if err != nil {
		return
	}

	val, ok := mp["token"]
	if ok {
		ZentaoToken = val.(string)
	}

	return
}
