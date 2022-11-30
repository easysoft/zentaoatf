package main

import (
	"flag"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	"github.com/easysoft/zentaoatf/internal/server/core/web"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"syscall"
)

var (
	AppVersion string
	BuildTime  string
	GoVersion  string
	GitHash    string

	uuid    = ""
	flagSet *flag.FlagSet
)

// @title ZTF服务端API文档
// @version 1.0
// @contact.name API Support
// @contact.url https://github.com/easysoft/zentaoatf/issues
// @contact.email 462626@qq.com
func main() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		cleanup()
		os.Exit(0)
	}()

	flagSet = flag.NewFlagSet("ztf", flag.ContinueOnError)

	flagSet.StringVar(&uuid, "uuid", "", "区分服务进程的唯一ID")

	flagSet.StringVar(&serverConfig.CONFIG.Server, "s", "http://pms.deeptest.loc", "")
	flagSet.StringVar(&serverConfig.CONFIG.Ip, "i", commConsts.Ip, "服务机器IP")
	flagSet.IntVar(&serverConfig.CONFIG.Port, "p", commConsts.Port, "服务端口")
	flagSet.StringVar(&serverConfig.CONFIG.Secret, "secret", "", "禅道认证安全码")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "run")
	}
	flagSet.Parse(os.Args[1:])

	switch os.Args[1] {
	case "version", "--version":
		logUtils.PrintVersion(AppVersion, BuildTime, GoVersion, GitHash)
		return
	default:
		commConsts.ExecFrom = commConsts.FromClient

		serverConfig.CONFIG.Server = httpUtils.AddSepIfNeeded(serverConfig.CONFIG.Server)

		webServer := web.Init(serverConfig.CONFIG.Port)
		if webServer == nil {
			return
		}

		websocketHelper.InitMq()

		webServer.Run()
	}

}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
