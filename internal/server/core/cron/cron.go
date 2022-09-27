package cron

import (
	"fmt"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	"github.com/easysoft/zentaoatf/internal/server/core/cache"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/easysoft/zentaoatf/pkg/lib/cron"
	"github.com/easysoft/zentaoatf/pkg/lib/date"
	"github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/kataras/iris/v12"
	"time"
)

type ServerCron struct {
	ExecService *service.ExecService `inject:""`
	JobService  *service.JobService  `inject:""`
}

func NewServerCron() *ServerCron {
	inst := &ServerCron{}
	return inst
}

func (s *ServerCron) Init() {
	cache.SyncMap.Store(cache.IsRunning, false)
	cache.SyncMap.Store(cache.LastLoopEndTime, int64(0))

	cronUtils.AddTask(
		"checkJob", fmt.Sprintf("@every %ds", serverConfig.JobCheckInterval),
		func() {
			isRunning, _ := cache.SyncMap.Load(cache.IsRunning)
			lastCompletedTime, _ := cache.SyncMap.Load(cache.LastLoopEndTime)

			if isRunning.(bool) || time.Now().Unix()-lastCompletedTime.(int64) < serverConfig.JobCheckInterval {
				logUtils.Infof("skip iteration" + dateUtils.DateTimeStr(time.Now()))
				return
			}

			cache.SyncMap.Store(isRunning, true)

			// start
			currJob := cache.GetCurrJob()
			s.JobService.Check(currJob)
			// end

			cache.SyncMap.Store(isRunning, false)
			cache.SyncMap.Store(lastCompletedTime, time.Now().Unix())
		},
	)

	iris.RegisterOnInterrupt(func() {
		cronUtils.Stop()
	})
}
