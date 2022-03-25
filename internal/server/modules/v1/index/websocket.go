package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type WebSocketModule struct {
	WebSocketCtrl *controller.WebSocketCtrl `inject:""`
}

func NewWebSocketModule() *WebSocketModule {
	return &WebSocketModule{}
}

// Party 项目
func (m *WebSocketModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())
	}
	return module.NewModule("/websocket", handler)
}
