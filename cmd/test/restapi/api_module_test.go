package main

import (
	"fmt"
	"testing"

	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	httpHelper "github.com/easysoft/zentaoatf/cmd/test/helper/http"
	"github.com/easysoft/zentaoatf/cmd/test/restapi/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/tidwall/gjson"
)

func TestModuleApi(t *testing.T) {
	suite.RunSuite(t, new(ModuleApiSuite))
}

type ModuleApiSuite struct {
	suite.Suite
}

func (s *ModuleApiSuite) BeforeEach(t provider.T) {
	commonTestHelper.ReplaceLabel(t, "ModuleApi")
}

func (s *ModuleApiSuite) TestModuleListForCaseApi(t provider.T) {
	t.ID("7636")
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("modules?type=case&id=%d", config.ProductId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	firstModuleId := gjson.Get(string(bodyBytes), "modules.0.id").Int()

	t.Require().Greater(firstModuleId, int64(0), "list modules failed")
}

func getModuleMinId() (id int64) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("modules?type=case&id=%d", config.ProductId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	modules := gjson.Get(string(bodyBytes), "modules").Array()
	for _, item := range modules {
		itemId := item.Get("id").Int()
		if id == 0 || (itemId > 0 && id > itemId) {
			id = itemId
		}
	}

	return
}
