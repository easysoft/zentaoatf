package service

import (
	"fmt"
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	cronUtils "github.com/easysoft/zentaoatf/src/server/utils/cron"
)

type CronService struct {
	commonService *CommonService
}

func NewCronService(commonService *CommonService) *CronService {
	return &CronService{commonService: commonService}
}

func (s *CronService) Init() {
	cronUtils.AddTaskFuc(
		"HeartBeat",

		fmt.Sprintf("@every %ds", serverConst.HeartBeatInterval),
		func() {
			s.commonService.HeartBeat()
		},
	)
}
