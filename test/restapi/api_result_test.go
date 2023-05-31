package main

import (
	"encoding/json"
	"fmt"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	httpHelper "github.com/easysoft/zentaoatf/test/helper/http"
	"github.com/easysoft/zentaoatf/test/restapi/config"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/tidwall/gjson"
	"log"
	"testing"
)

func TestResultApi(t *testing.T) {
	suite.RunSuite(t, new(ResultApiSuite))
}

type ResultApiSuite struct {
	suite.Suite
}

func (s *ResultApiSuite) BeforeEach(t provider.T) {
	t.AddSubSuite("SuiteApi")
}

//func (s *ResultApiSuite) TestResultSubmitZtfResultApi(t provider.T) {
//	t.ID("7626,7628")
//	token := httpHelper.Login()
//
//	latestId := getCaseResult(config.CaseId)["latestId"].(int64)
//
//	url := zentaoHelper.GenApiUrl("ciresults", nil, constTestHelper.ZentaoSiteUrl)
//
//	report := commDomain.ZtfReport{}
//	err := json.Unmarshal([]byte(ztfReportJson), &report)
//	t.Require().Equal(err, nil, "submit result failed")
//
//	report.Name = "接口测试任务" + stringUtils.NewUuid()
//
//	resp, err := httpHelper.Post(url, token, report)
//	log.Print(resp)
//	t.Require().Equal(err, nil, "submit result failed")
//
//	// check case result
//	latestIdNew := getCaseResult(config.CaseId)["latestId"].(int64)
//	t.Require().Equal(latestIdNew, latestId+1, "submit result failed")
//
//	// check task record
//	tasksBytes := listTask(token)
//	name := gjson.Get(string(tasksBytes), "testtasks.0.name").String()
//	t.Require().Equal(name, report.Name, "submit result failed")
//}

func (s *ResultApiSuite) TestResultSubmitUnitResultApi(t provider.T) {
	t.ID("7627")
	token := httpHelper.Login()

	latestCaseResultId := getCaseResult(config.CaseId)["latestId"].(int64)

	url := zentaoHelper.GenApiUrl("ciresults", nil, constTestHelper.ZentaoSiteUrl)

	report := commDomain.ZtfReport{}
	err := json.Unmarshal([]byte(unitReportJson), &report)
	t.Require().Equal(err, nil, "submit result failed")
	report.Name = "接口测试任务" + stringUtils.NewUuid()

	resp, err := httpHelper.Post(url, token, report)
	log.Print(resp)
	t.Require().Equal(err, nil, "submit result failed")

	// check case result
	latestCaseResultIdNew := getCaseResult(config.CaseId)["latestId"].(int64)
	t.Require().Greater(latestCaseResultIdNew, latestCaseResultId, "submit result failed")

	// check task record
	//tasksBytes := listTask(token)
	//name := gjson.Get(string(tasksBytes), "testtasks.0.name").String()
	//t.Require().Equal(name, report.Name, "submit result failed")
}

//func (s *ResultApiSuite) TestResultSubmitSameTaskIdApi(t provider.T) {
//	t.ID("7630")
//	token := httpHelper.Login()
//
//	latestCaseResultId := getCaseResult(config.CaseId)["latestId"]
//
//	url := zentaoHelper.GenApiUrl("ciresults", nil, constTestHelper.ZentaoSiteUrl)
//
//	report := commDomain.ZtfReport{}
//	err := json.Unmarshal([]byte(ztfReportJson), &report)
//	t.Require().Equal(err, nil, "submit result failed")
//
//	_, err = httpHelper.Post(url, token, report)
//	t.Require().Equal(err, nil, "submit result failed")
//
//	// check case result
//	latestCaseResultId2 := getCaseResult(config.CaseId)["latestId"]
//	t.Require().Greater(latestCaseResultId2, latestCaseResultId, "submit result failed")
//
//	// get latest task id
//	latestTaskId := getLatestTaskId(token)
//
//	// submit again with same task id
//	report.TaskId = latestTaskId
//	_, err = httpHelper.Post(url, token, report)
//	t.Require().Equal(err, nil, "submit result failed")
//
//	// check case result
//	latestCaseResultId3 := getCaseResult(config.CaseId)["latestId"].(int64)
//	t.Require().Greater(latestCaseResultId3, latestCaseResultId2, "submit result failed")
//
//	// get latest task id
//	latestTaskId2 := getLatestTaskId(token)
//
//	// check not add an new task
//	t.Require().Equal(latestTaskId2, latestTaskId, "submit result failed")
//}

func getCaseResult(caseId int) (result map[string]interface{}) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testcases/%d/results", caseId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	result = map[string]interface{}{}

	result["latestId"] = gjson.Get(string(bodyBytes), "results.0.id").Int()

	return
}

const ztfReportJson = `
{
	"name": "",
	"platform": "mac",
	"testType": "func",
	"testTool": "ztf",
	"buildTool": "",
	"testCommand": "/private/var/folders/ry/yjxnkwt12kl6d1d13q5cz6wc0000gn/T/GoLand/___test_1 run demo/t/test.py -p 1",
	"workspaceType": "",
	"workspacePath": "/Users/aaron/rd/project/zentao/go/ztf/",
	"submitResult": false,
	"execBy": "case",
	"productId": 1,
	"zentaoData": "",
	"buildUrl": "",
	"log": "2023-05-29 13:48:35.118\texec/script.go:35\t开始执行/Users/aaron/rd/project/zentao/go/ztf/demo/t/test.py于2023-05-29 13:48:35。\n2023-05-29 13:48:35.521\texec/file.go:117\tone\n\n2023-05-29 13:48:35.521\texec/file.go:117\ttwo\n\n2023-05-29 13:48:35.521\texec/file.go:117\tthree\n\n2023-05-29 13:48:35.537\texec/script.go:64\t结束执行/Users/aaron/rd/project/zentao/go/ztf/demo/t/test.py于2023-05-29 13:48:35。",
	"pass": 1,
	"fail": 0,
	"skip": 0,
	"total": 1,
	"startTime": 1685339315,
	"endTime": 1685339315,
	"duration": 0,
	"funcResult": [
		{
			"id": 1,
			"workspaceId": 0,
			"seq": "",
			"key": "e515b2623360a567023c8ddd9e91514c",
			"productId": 1,
			"path": "/Users/aaron/rd/project/zentao/go/ztf/demo/t/test.py",
			"relativePath": "demo/t/test.py",
			"status": "pass",
			"title": "测试返回结果",
			"steps": [
				{
					"id": "1",
					"name": "返回 one",
					"status": "pass",
					"checkPoints": [
						{
							"numb": 1,
							"expect": "one",
							"actual": "one",
							"status": "pass"
						}
					]
				},
				{
					"id": "2",
					"name": "返回 two",
					"status": "pass",
					"checkPoints": [
						{
							"numb": 1,
							"expect": "two",
							"actual": "two",
							"status": "pass"
						}
					]
				},
				{
					"id": "3",
					"name": "返回 three",
					"status": "pass",
					"checkPoints": [
						{
							"numb": 1,
							"expect": "three",
							"actual": "three",
							"status": "pass"
						}
					]
				}
			]
		}
	]
}
`

const unitReportJson = `
{
    "name": "",
    "platform": "mac",
    "testType": "unit",
    "testTool": "gotest",
    "buildTool": "",
    "testCommand": "go test restapi/api_product_test.go -v",
    "workspaceType": "",
    "workspacePath": "",
    "submitResult": false,
	"productId": 1,
    "zentaoData": "",
    "buildUrl": "",
    "log": "2023-05-31 13:58:03.483\texec/unit.go:155\t=== RUN ...",
    "pass": 2,
    "fail": 0,
    "skip": 0,
    "total": 2,
    "startTime": 1685512678,
    "endTime": 1685512683,
    "duration": 5,
    "unitResult":
    [
        {
            "title": "TestProductDetailApi",
            "testSuite": "ProductApiSuite-ProductApi",
			"productId": 1,
            "startTime": 0,
            "endTime": 0,
            "duration": 0.081,
            "failure": null,
            "errorType": "",
            "errorContent": "",
            "id": 1,
            "cid": 1,
            "status": "pass"
        },
        {
            "title": "TestProductListApi",
            "testSuite": "ProductApiSuite-ProductApi",
            "startTime": 0,
            "endTime": 0,
            "duration": 0.165,
            "failure": null,
            "errorType": "",
            "errorContent": "",
            "id": 2,
            "cid": 2,
            "status": "pass"
        }
    ]
}

`
