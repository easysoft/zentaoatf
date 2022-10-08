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

		index.Post("/add", m.JobCtrl.Add).Name = "添加任务到队列"
		index.Post("/remove", m.JobCtrl.Remove).Name = "移除队列中任务"
		index.Post("/stop", m.JobCtrl.Stop).Name = "终止当前执行的任务"
	}
	return module.NewModule("/jobs", handler)
}
