package main

import (
	"flag"
	"fmt"
	"testing"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	uiTest "github.com/easysoft/zentaoatf/test/helper/zentao/ui"
)

func main() {
	commConsts.ExecFrom = commConsts.FromCmd
	serverConfig.InitLog()
	serverConfig.InitExecLog(commonTestHelper.RootPath)
	commConsts.ZtfDir = commonTestHelper.RootPath
	i118Utils.Init("zh-CN", commConsts.AppServer)
	// commConsts.Verbose = true
	var version = flag.String("zentaoVersion", "", "")
	testing.Init()
	flag.Parse()
	fmt.Println(*version)
	defer func() {
		execHelper.KillProcessByUUID("ui_auto_test")
		uiTest.Close()
	}()

	err := commonTestHelper.InitZentao(*version)
	if err != nil {
		fmt.Println("Init zentao data fail ", err)
	}
	// err = commonTestHelper.Pull()
	// if err != nil {
	// 	fmt.Println("Git pull code fail ", err)
	// }
	err = commonTestHelper.BuildCli()
	if err != nil {
		fmt.Println("Build cli fail ", err)
	}
	err = commonTestHelper.RunServer()
	if err != nil {
		fmt.Println("Build server fail ")
	}
	err = commonTestHelper.RunUi()
	if err != nil {
		fmt.Println("Build server fail ")
	}
	commonTestHelper.TestUi()
}
