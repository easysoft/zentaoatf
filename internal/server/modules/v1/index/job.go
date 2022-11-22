package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type JobModule struct {
	JobCtrl *controller.JobCtrl `inject:""`
}

func NewJobModule() *JobModule {
	return &JobModule{}
}

// Party 执行
func (m *JobModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		//index.Use(core.Auth())

		index.Post("/add", m.JobCtrl.Add).Name = "添加任务到队列"
		index.Post("/cancel", m.JobCtrl.Cancel).Name = "强制终止任务"
	}
	return module.NewModule("/jobs", handler)
}
