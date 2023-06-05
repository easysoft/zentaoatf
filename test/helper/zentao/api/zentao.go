package apiTest

import (
	"fmt"

	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	httpHelper "github.com/easysoft/zentaoatf/test/helper/http"
	"github.com/easysoft/zentaoatf/test/restapi/config"
	"github.com/tidwall/gjson"
)

func GetCaseTitleById(id int) string {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testcases/%d", id), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	return gjson.Get(string(bodyBytes), "title").String()
}

func GetCaseResult(caseId int) (result map[string]interface{}) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testcases/%d/results", caseId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	result = map[string]interface{}{}

	result["Id"] = gjson.Get(string(bodyBytes), "results.0.id").Int()
	result["CaseResult"] = gjson.Get(string(bodyBytes), "results.0.caseResult").String()
	result["Date"] = gjson.Get(string(bodyBytes), "results.0.date").String()

	return
}

func GetLastBugId() int64 {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("products/%d/bugs", config.ProductId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	return gjson.Get(string(bodyBytes), "bugs.0.id").Int()
}
