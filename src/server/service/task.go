package service

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/server/domain"
	serverUtils "github.com/easysoft/zentaoatf/src/server/utils/common"
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"strconv"
	"time"
)

var (
	tasks     = make([]domain.Build, 0)
	tm        = time.Now()
	isRunning = false
)

type TaskService struct {
	buildService *BuildService
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (s *TaskService) Add(build domain.Build) {
	tasks = append(tasks, build)
}

func (s *TaskService) Peek() domain.Build {
	return tasks[0]
}

func (s *TaskService) Remove() (task domain.Build) {
	if len(tasks) == 0 {
		return task
	}

	task = tasks[0]
	tasks = tasks[1:]

	return task
}

func (s *TaskService) Start() {
	tm = time.Now()
	isRunning = true
}
func (s *TaskService) End() {
	isRunning = false
}

func (s *TaskService) GetSize() int {
	return len(tasks)
}

func (s *TaskService) CheckRunning() bool {
	if time.Now().Unix()-tm.Unix() > serverConst.AgentRunTime*60*1000 {
		isRunning = false
	}
	return isRunning
}

func (s *TaskService) ListTask() (data []domain.Build) {
	data = tasks
	return
}

func (s *TaskService) ListHistory() (data []map[string]string) {
	data = serverUtils.ListHistoryLog()

	for key, item := range data {
		data[key]["url"] = fmt.Sprintf("http://%s:%s/down?f=%s", vari.IP, strconv.Itoa(vari.Port), item["name"])
	}

	return
}
