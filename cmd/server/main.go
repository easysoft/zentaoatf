package main

import (
	"flag"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/core/cron"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	"github.com/aaronchen2k/deeptest/internal/server/core/web"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1"
	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"
)

var (
	port = 0
	uuid = ""
)

// @title ZTF服务端API文档
// @version 1.0
// @contact.name API Support
// @contact.url https://github.com/easysoft/zentaoatf/issues
// @contact.email 462626@qq.com
func main() {
	port = *flag.Int("p", 0, "服务端口")
	uuid = *flag.String("uuid", "", "区分服务进程的唯一ID")
	flag.Parse()

	webServer := web.Init(port)
	if webServer == nil {
		return
	}

	cron.NewServerCron().Init()

	injectModule(webServer)

	webServer.Run()
}

func injectModule(ws *web.WebServer) {
	var g inject.Graph
	g.Logger = logrus.StandardLogger()

	indexModule := v1.NewIndexModule()

	// inject objects
	if err := g.Provide(
		&inject.Object{Value: dao.GetDB()},
		&inject.Object{Value: indexModule},
	); err != nil {
		logrus.Fatalf("provide usecase objects to the Graph: %v", err)
	}
	err := g.Populate()
	if err != nil {
		logrus.Fatalf("populate the incomplete Objects: %v", err)
	}

	ws.AddModule(indexModule.Party())

	logUtils.Infof("start server")
}
