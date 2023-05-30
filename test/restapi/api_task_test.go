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

func TestTaskApi(t *testing.T) {
	suite.RunSuite(t, new(TaskApiSuite))
}

type TaskApiSuite struct {
	suite.Suite
}

func (s *TaskApiSuite) BeforeEach(t provider.T) {
	t.ID("0")
	t.AddSubSuite("TaskApi")
}

func (s *TaskApiSuite) TestTaskListApi(t provider.T) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testtasks?product=%d", ProductId), nil, constTestHelper.ZtfUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	firstTaskId := gjson.Get(string(bodyBytes), "testtasks.0.id").Int()

	t.Require().Greater(firstTaskId, int64(0), "list testtasks failed")
}

func (s *TaskApiSuite) TestTaskDetailApi(t provider.T) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testtasks/%d", TaskId), nil, constTestHelper.ZtfUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	name := gjson.Get(string(bodyBytes), "name").String()

	t.Require().Greater(len(name), 0, "get testsuite failed")
}
