package main

import (
	"flag"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	"github.com/easysoft/zentaoatf/internal/server/core/web"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"os"
)

var (
	AppVersion string
	BuildTime  string
	GoVersion  string
	GitHash    string

	uuid = ""
)

// @title ZTF服务端API文档
// @version 1.0
// @contact.name API Support
// @contact.url https://github.com/easysoft/zentaoatf/issues
// @contact.email 462626@qq.com
func main() {
	flag.StringVar(&uuid, "uuid", "", "区分服务进程的唯一ID")

	flag.StringVar(&serverConfig.CONFIG.Server, "s", "http://pms.deeptest.loc", "")
	flag.StringVar(&serverConfig.CONFIG.Ip, "i", commConsts.Ip, "服务所在机器IP地址")
	flag.IntVar(&serverConfig.CONFIG.Port, "p", commConsts.Port, "服务所在端口")
	flag.StringVar(&serverConfig.CONFIG.Secret, "secret", "", "禅道认证安全码")

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

		webServer := web.Init(serverConfig.CONFIG.Port)
		if webServer == nil {
			return
		}

		websocketHelper.InitMq()

		webServer.Run()
	}

}
