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
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	uiTest "github.com/easysoft/zentaoatf/test/helper/zentao/ui"
)

func main() {
	commConsts.ExecFrom = commConsts.FromCmd
	commConsts.Verbose = true
	serverConfig.InitLog()
	serverConfig.InitExecLog(constTestHelper.RootPath)
	commConsts.ZtfDir = constTestHelper.RootPath
	i118Utils.Init("zh-CN", commConsts.AppServer)
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
	// err = constTestHelper.Pull()
	// if err != nil {
	// 	fmt.Println("Git pull code fail ", err)
	// }
	err = commonTestHelper.BuildCli()
	if err != nil {
		fmt.Println("Build cli fail ", err)
	}
	commonTestHelper.TestCli()
}
