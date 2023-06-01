package main

import (
	"fmt"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	httpHelper "github.com/easysoft/zentaoatf/test/helper/http"
	"github.com/easysoft/zentaoatf/test/restapi/config"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/tidwall/gjson"
	"testing"
)

func TestSuiteApi(t *testing.T) {
	suite.RunSuite(t, new(SuiteApiSuite))
}

type SuiteApiSuite struct {
	suite.Suite
}

func (s *SuiteApiSuite) BeforeEach(t provider.T) {
	t.AddSubSuite("SuiteApi")
}

func (s *SuiteApiSuite) TestSuiteListApi(t provider.T) {
	t.ID("7622")
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("products/%d/testsuites", config.ProductId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	firstSuiteId := gjson.Get(string(bodyBytes), "testsuites.0.id").Int()

	t.Require().Greater(firstSuiteId, int64(0), "list testsuite failed")
}

func (s *SuiteApiSuite) TestSuiteDetailApi(t provider.T) {
	t.ID("7623")
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testsuites/%d", config.SuiteId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	name := gjson.Get(string(bodyBytes), "name").String()

	t.Require().Greater(len(name), 0, "get testsuite failed")
}

func getSuiteMinId() (id int64) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("products/%d/testsuites", config.ProductId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	suites := gjson.Get(string(bodyBytes), "testsuites").Array()
	for _, suite := range suites {
		suiteId := suite.Get("id").Int()
		if id == 0 || (suiteId > 0 && id > suiteId) {
			id = suiteId
		}
	}

	return
}
