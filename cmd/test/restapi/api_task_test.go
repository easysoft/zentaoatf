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

func TestTaskApi(t *testing.T) {
	suite.RunSuite(t, new(TaskApiSuite))
}

type TaskApiSuite struct {
	suite.Suite
}

func (s *TaskApiSuite) BeforeEach(t provider.T) {
	commonTestHelper.ReplaceLabel(t, "TaskApi")
}

func (s *TaskApiSuite) TestTaskListApi(t provider.T) {
	t.ID("7624")
	token := httpHelper.Login()

	bodyBytes := listTask(token)

	firstTaskId := gjson.Get(string(bodyBytes), "testtasks.0.id").Int()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testtasks?product=%d", config.ProductId), nil, constTestHelper.ZentaoSiteUrl)

	t.Require().Greater(firstTaskId, int64(0), "list testtasks failed, url: "+url)
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

func getTaskMinId() (id int64) {
	token := httpHelper.Login()

	url := zentaoHelper.GenApiUrl(fmt.Sprintf("testtasks?product=%d", config.ProductId), nil, constTestHelper.ZentaoSiteUrl)

	bodyBytes, _ := httpHelper.Get(url, token)

	tasks := gjson.Get(string(bodyBytes), "testtasks").Array()
	for _, task := range tasks {
		taskId := task.Get("id").Int()
		if id == 0 || (taskId > 0 && id > taskId) {
			id = taskId
		}
	}

	return
}
