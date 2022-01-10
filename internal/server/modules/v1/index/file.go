package index

import (
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type FileModule struct {
	FileCtrl *controller.FileCtrl `inject:""`
}

func NewFileModule() *FileModule {
	return &FileModule{}
}

// Party 上传文件模块
func (m *FileModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck())
		index.Post("/upload", iris.LimitRequestBodySize(serverConfig.CONFIG.MaxSize+1<<20), m.FileCtrl.Upload).Name = "上传文件"
		index.Get("/listDir", m.FileCtrl.ListDir).Name = "获取目录文件"
	}
	return module.NewModule("/file", handler)
}
