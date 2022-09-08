package index

import (
	"github.com/easysoft/zentaoatf/internal/server/core/module"
	"github.com/easysoft/zentaoatf/internal/server/middleware"
	"github.com/kataras/iris/v12"
)

type HeartBeatModule struct {
}

func NewHeartBeatModule() *HeartBeatModule {
	return &HeartBeatModule{}
}

// Party 执行
func (m *HeartBeatModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())

		index.Get("/", func(ctx iris.Context) {
			ctx.JSON(iris.Map{"code": 0, "msg": "ok"})
		})
	}
	return module.NewModule("/heartbeat", handler)
}
