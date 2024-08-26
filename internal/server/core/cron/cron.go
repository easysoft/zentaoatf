package cron

import (
	"fmt"

	"github.com/kataras/iris/v12"

	"github.com/easysoft/zentaoatf/internal/server/core/cache"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"

	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	cronUtils "github.com/easysoft/zentaoatf/pkg/lib/cron"
)

type ServerCron struct {
	HeartbeatService *service.HeartbeatService `inject:""`
	JobService       *service.JobService       `inject:""`
}

func NewServerCron() *ServerCron {
	inst := &ServerCron{}
	return inst
}

func (s *ServerCron) Init() {
	cache.SyncMap.Store(cache.IsRunning, false)
	cache.SyncMap.Store(cache.LastLoopEndTime, int64(0))

	if serverConfig.CONFIG.Host != "" {
		cronUtils.AddTask(
			"heartbeat", fmt.Sprintf("@every %ds", serverConfig.HeartbeatInterval),
			func() {
				s.HeartbeatService.Heartbeat()
			},
		)
	}

	cronUtils.AddTask(
		"checkJob", fmt.Sprintf("@every %ds", serverConfig.JobCheckInterval),
		func() {
			s.JobService.Check()
		},
	)

	iris.RegisterOnInterrupt(func() {
		cronUtils.Stop()
	})
}
