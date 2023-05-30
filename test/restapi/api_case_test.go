package main

import (
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	httpHelper "github.com/easysoft/zentaoatf/test/helper/http"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/tidwall/gjson"
	"testing"
)

func TestCaseApi(t *testing.T) {
	suite.RunSuite(t, new(CaseApiSuite))
}

type CaseApiSuite struct {
	suite.Suite
}

func (s *CaseApiSuite) BeforeEach(t provider.T) {
	t.ID("0")
	t.AddSubSuite("CaseApi")
}

func (s *CaseApiSuite) TestCaseListApi(t provider.T) {
	token := httpHelper.Login()

	params := map[string]interface{}{
		"limit": 10,
	}
	url := zentaoHelper.GenApiUrl(fmt.Sprintf("/products/%d/testcases", ProductId), params, constTestHelper.ZtfUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	firstCaseId := gjson.Get(string(bodyBytes), "testcases.0.id").Int()

	t.Require().Greater(firstCaseId, int64(0), "list testcases failed")
}

func (s *CaseApiSuite) TestCaseListByModuleApi(t provider.T) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("/products/%d/testcases?module=%d", ProductId, ModuleId),
		nil, constTestHelper.ZtfUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	firstCaseId := gjson.Get(string(bodyBytes), "testcases.0.id").Int()

	t.Require().Greater(firstCaseId, int64(0), "list testcases failed")
}

func (s *CaseApiSuite) TestCaseListBySuiteApi(t provider.T) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("/testsuites/%d", SuiteId), nil, constTestHelper.ZtfUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	firstCaseId := gjson.Get(string(bodyBytes), "testcases.0.id").Int()

	t.Require().Greater(firstCaseId, int64(0), "list testcases failed")
}

func (s *CaseApiSuite) TestCaseListByTaskApi(t provider.T) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("/testtasks/%d", TaskId), nil, constTestHelper.ZtfUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	firstCaseId := gjson.Get(string(bodyBytes), "testcases.0.id").Int()

	t.Require().Greater(firstCaseId, int64(0), "list testcases failed")
}

func (s *CaseApiSuite) TestCaseDetailApi(t provider.T) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testcases/%d", CaseId), nil, constTestHelper.ZtfUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	title := gjson.Get(string(bodyBytes), "title").String()

	t.Require().Greater(len(title), 0, "get testcases failed")
}

func (s *CaseApiSuite) TestCaseCheckinApi(t provider.T) {
	token := httpHelper.Login()
	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testcases/%d", CaseId), nil, constTestHelper.ZtfUrl)

	steps := []commDomain.ZentaoCaseStep{
		{Type: commConsts.Item, Desc: "Step 1", Expect: "Expect 1"},
		{Type: commConsts.Item, Desc: "Step 2", Expect: "Expect 1"},
		{Type: commConsts.Item, Desc: "Step 3", Expect: "Expect 1"},
	}

	title := "用例新名字" + stringUtils.NewUuid()
	requestObj := map[string]interface{}{
		"type":  "feature",
		"title": title,
		"steps": steps,

		"path":   "path_of_case",
		"script": "script_of_case",
		"lang":   "php",
	}

	bodyBytes, _ := httpHelper.Put(url, token, requestObj)

	actualTitle := gjson.Get(string(bodyBytes), "title").String()
	t.Require().Equal(actualTitle, title, "checkin testcases failed")

	newCase := getCase(CaseId)
	titleFromRemote := newCase["title"]
	t.Require().Equal(titleFromRemote, title, "get testcases failed")
}

func getCase(id int) (cs map[string]interface{}) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testcases/%d", id), nil, constTestHelper.ZtfUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	cs = map[string]interface{}{}

	cs["id"] = gjson.Get(string(bodyBytes), "id").Int()
	cs["title"] = gjson.Get(string(bodyBytes), "title").String()

	return
}
