package service

import (
	"fmt"
	serverUtils "github.com/easysoft/zentaoatf/src/server/utils/common"
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	cronUtils "github.com/easysoft/zentaoatf/src/server/utils/cron"
)

type CronService struct {
	heartBeatService *HeartBeatService

	buildService *BuildService
	taskService  *TaskService
	execService  *ExecService
}

func NewCronService(heartBeatService *HeartBeatService,
	buildService *BuildService, taskService *TaskService,
	execService *ExecService) *CronService {
	return &CronService{heartBeatService: heartBeatService,
		buildService: buildService, taskService: taskService, execService: execService}
}

func (s *CronService) Init() {
	cronUtils.AddTaskFuc(
		"HeartBeat",

		fmt.Sprintf("@every %ds", serverConst.HeartBeatInterval),
		func() {
			if serverUtils.IsVmAgent() { // vm
				// is running，register busy
				if s.taskService.CheckRunning() {
					s.heartBeatService.HeartBeat(true)
					return
				}

				// no task to run, submit free
				if s.taskService.GetSize() == 0 {
					s.heartBeatService.HeartBeat(false)
					return
				}

				// has task to run，register busy, then run
				build := s.taskService.Peek()
				s.heartBeatService.HeartBeat(true)
				s.execService.Exec(build)
			}
		},
	)
}
