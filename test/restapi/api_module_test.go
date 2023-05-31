package main

import (
	"fmt"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	httpHelper "github.com/easysoft/zentaoatf/test/helper/http"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/tidwall/gjson"
	"testing"
)

func TestModuleApi(t *testing.T) {
	suite.RunSuite(t, new(ModuleApiSuite))
}

type ModuleApiSuite struct {
	suite.Suite
}

func (s *ModuleApiSuite) BeforeEach(t provider.T) {
	t.AddSubSuite("ModuleApi")
}

func (s *ModuleApiSuite) TestModuleListForCaseApi(t provider.T) {
	t.ID("7636")
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("modules?type=case&id=%d", ProductId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	firstModuleId := gjson.Get(string(bodyBytes), "modules.0.id").Int()

	t.Require().Greater(firstModuleId, int64(0), "list modules failed")
}
