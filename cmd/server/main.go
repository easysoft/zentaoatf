package main

import (
	"flag"
	"github.com/aaronchen2k/deeptest/internal/server/core/cron"
	"github.com/aaronchen2k/deeptest/internal/server/core/web"
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
	flag.IntVar(&port, "p", 0, "服务端口")
	flag.StringVar(&uuid, "uuid", "", "区分服务进程的唯一ID")
	flag.Parse()

	webServer := web.Init(port)
	if webServer == nil {
		return
	}

	cron.NewServerCron().Init()

	webServer.Run()
}
