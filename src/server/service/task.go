package service

import (
	"github.com/easysoft/zentaoatf/src/server/domain"
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
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
