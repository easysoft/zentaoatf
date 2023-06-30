package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	uiTest "github.com/easysoft/zentaoatf/cmd/test/helper/zentao/ui"
	commandConfig "github.com/easysoft/zentaoatf/internal/command/config"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
)

var (
	runFrom, version, testToRun string
	flagSet                     *flag.FlagSet
)

func main() {
	defer func() {
		execHelper.KillProcessByUUID("ui_auto_test")
		uiTest.Close()
	}()

	flagSet = flag.NewFlagSet("restapi", flag.ContinueOnError)

	flagSet.StringVar(&runFrom, "runFrom", "cmd", "")
	flagSet.StringVar(&runFrom, "f", "cmd", "")

	flagSet.StringVar(&version, "version", "latest", "")
	flagSet.StringVar(&version, "v", "latest", "")

	flagSet.StringVar(&testToRun, "testToRun", "", "")
	flagSet.StringVar(&testToRun, "t", "", "")

	testing.Init()
	flagSet.Parse(os.Args[1:])
	fmt.Println(version)

	commConsts.ExecFrom = commConsts.FromCmd
	commConsts.ZtfDir = constTestHelper.RootPath

	commandConfig.InitLog()
	serverConfig.InitExecLog(constTestHelper.RootPath)
	i118Utils.Init("zh-CN", commConsts.AppServer)

	if runFrom == "jenkins" {
		constTestHelper.ZentaoSiteUrl = constTestHelper.ZentaoSiteUrl[:strings.LastIndex(constTestHelper.ZentaoSiteUrl, ":")]

		err := commonTestHelper.InitZentaoData()
		if err != nil {
			fmt.Println("Init zentao data fail ", err)
		}
	} else {
		err := commonTestHelper.InitZentao(version)
		if err != nil {
			fmt.Println("Init zentao data fail ", err)
		}

		fmt.Println("build cli")
		err = commonTestHelper.BuildCli()
		if err != nil {
			fmt.Println("Build cli fail ", err)
		}

		fmt.Println("run server")
		err = commonTestHelper.RunServer()
		if err != nil {
			fmt.Println("Build server fail ")
		}

		fmt.Println("run ui")
		err = commonTestHelper.RunUi()
		if err != nil {
			fmt.Println("Build server fail ")
		}
	}

	commonTestHelper.WaitZtfAccessed()
	commonTestHelper.TestUi(version)
}
