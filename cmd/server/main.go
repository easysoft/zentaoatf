package main

import (
	"flag"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/core/cron"
	"github.com/aaronchen2k/deeptest/internal/server/core/web"
	"os"
)

var (
	appVersion string
	buildTime  string
	goVersion  string
	gitHash    string

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

	switch os.Args[1] {
	case "version", "--version":
		logUtils.PrintVersion(appVersion, buildTime, goVersion, gitHash)
	default:
	}

	webServer := web.Init(port)
	if webServer == nil {
		return
	}

	cron.NewServerCron().Init()

	webServer.Run()
}
