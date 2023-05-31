package main

import (
	"fmt"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	httpHelper "github.com/easysoft/zentaoatf/test/helper/http"
	"github.com/easysoft/zentaoatf/test/restapi/config"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/tidwall/gjson"
	"testing"
)

const (
	BugTitle = "缺陷提交测试"
	BugStep  = `步骤1： step1 fail<br />  验证点：fail<br />    期待结果：<br />      expect1<br />    实际结果：<br />      3<br /><br />步骤2： step2 fail<br />  验证点：fail<br />    期待结果：<br />      expect2<br />    实际结果：<br />      2<br /><br />步骤3： step3 fail<br />  验证点：fail<br />    期待结果：<br />      expect3<br />    实际结果：<br />      1<br />`
)

func TestBugApi(t *testing.T) {
	suite.RunSuite(t, new(BugApiSuite))
}

type BugApiSuite struct {
	suite.Suite
}

func (s *BugApiSuite) BeforeEach(t provider.T) {
	t.ID("0")
	t.AddSubSuite("BugApi")
}

func (s *BugApiSuite) TestBugListApi(t provider.T) {
	t.ID("7617")
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("products/%d/bugs", config.ProductId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	firstBugId := gjson.Get(string(bodyBytes), "bugs.0.id").Int()

	t.Require().Greater(firstBugId, int64(0), "list bug failed")
}

func (s *BugApiSuite) TestBugCreateApi(t provider.T) {
	t.ID("7619")
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("products/%d/bugs", config.ProductId), nil, constTestHelper.ZentaoSiteUrl)

	title := BugTitle + stringUtils.NewUuid()

	data := map[string]interface{}{
		"id":       0,
		"title":    title,
		"type":     "codeerror",
		"product":  config.ProductId,
		"module":   0,
		"severity": 3,
		"pri":      3,
		"case":     0,
		"steps":    BugStep,
		"uid":      stringUtils.NewUuid(),
		"openedBuild": []string{
			"trunk",
		},
		"caseVersion": "0",
		"oldTaskID":   "0",
		"openedDate":  "",
		"openedBy":    "",
		"statusName":  "",
	}

	bodyBytes, _ := httpHelper.Post(url, token, data)

	newBugId := gjson.Get(string(bodyBytes), "id").Int()
	t.Require().Greater(newBugId, int64(0), "create bug failed")

	newBug := getBug(newBugId)
	newBugTitleFromRemote := newBug["title"]
	t.Require().Equal(newBugTitleFromRemote, title, "get new bug failed")
}

func (s *BugApiSuite) TestBugOptionsApi(t provider.T) {
	t.ID("7618")
	token := httpHelper.Login()

	params := map[string]interface{}{
		"product": config.ProductId,
	}
	url := zentaoHelper.GenApiUrl("options/bug", params, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	truckName := gjson.Get(string(bodyBytes), "options.build.trunk").String()

	t.Require().Equal(truckName, "主干", "list product failed")
}

func getBug(id int64) (bug map[string]interface{}) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("bugs/%d", id), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	bug = map[string]interface{}{}

	bug["id"] = gjson.Get(string(bodyBytes), "id").Int()
	bug["title"] = gjson.Get(string(bodyBytes), "title").String()

	return
}
