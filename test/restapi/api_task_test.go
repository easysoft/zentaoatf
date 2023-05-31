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

func TestTaskApi(t *testing.T) {
	suite.RunSuite(t, new(TaskApiSuite))
}

type TaskApiSuite struct {
	suite.Suite
}

func (s *TaskApiSuite) BeforeEach(t provider.T) {
	t.AddSubSuite("TaskApi")
}

func (s *TaskApiSuite) TestTaskListApi(t provider.T) {
	t.ID("7624")
	token := httpHelper.Login()

	bodyBytes := listTask(token)

	firstTaskId := gjson.Get(string(bodyBytes), "testtasks.0.id").Int()

	t.Require().Greater(firstTaskId, int64(0), "list testtasks failed")
}

func (s *TaskApiSuite) TestTaskDetailApi(t provider.T) {
	t.ID("7625")
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testtasks/%d", config.TaskId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	name := gjson.Get(string(bodyBytes), "name").String()

	t.Require().Greater(len(name), 0, "get testsuite failed")
}

func listTask(token string) (bodyBytes []byte) {
	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testtasks?product=%d", config.ProductId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ = httpHelper.Get(url, token)

	return
}

func getLatestTaskId(token string) (id int) {
	tasksBytes := listTask(token)
	idInt64 := gjson.Get(string(tasksBytes), "testtasks.0.id").Int()

	id = int(idInt64)

	return
}
