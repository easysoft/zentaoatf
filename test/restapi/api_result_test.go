package main

import (
	"encoding/json"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	httpHelper "github.com/easysoft/zentaoatf/test/helper/http"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestResultApi(t *testing.T) {
	suite.RunSuite(t, new(ResultApiSuite))
}

type ResultApiSuite struct {
	suite.Suite
}

func (s *ResultApiSuite) BeforeEach(t provider.T) {
	t.ID("0")
	t.AddSubSuite("SuiteApi")
}

func (s *ResultApiSuite) TestResultSubmitApi(t provider.T) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl("ciresults", nil, constTestHelper.ZentaoSiteUrl)

	report := commDomain.ZtfReport{}
	json.Unmarshal([]byte(reportJson), &report)

	_, err := httpHelper.Post(url, token, report)

	t.Require().Equal(err, nil, "list testsuite failed")

	//firstSuiteId := gjson.Get(string(bodyBytes), "testsuites.0.id").Int()
	//t.Require().Greater(firstSuiteId, int64(0), "list testsuite failed")
}

const reportJson = `
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
