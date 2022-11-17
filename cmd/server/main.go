package main

import (
	"flag"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	"github.com/easysoft/zentaoatf/internal/server/core/web"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"os"
)

var (
	AppVersion string
	BuildTime  string
	GoVersion  string
	GitHash    string

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

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "run")
	}

	switch os.Args[1] {
	case "version", "--version":
		logUtils.PrintVersion(AppVersion, BuildTime, GoVersion, GitHash)
		return
	default:
		commConsts.ExecFrom = commConsts.FromClient

		webServer := web.Init(port)
		if webServer == nil {
			return
		}

		websocketHelper.InitMq()

		webServer.Run()
	}

}
