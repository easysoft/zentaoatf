package main

import (
	"flag"
	"fmt"
	"testing"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
)

func main() {
	serverConfig.InitLog()
	commConsts.Verbose = true
	var version = flag.String("zentaoVersion", "", "")
	testing.Init()
	flag.Parse()
	fmt.Println(*version)
	err := commonTestHelper.InitZentao(*version)
	if err != nil {
		fmt.Println("Init zentao data fail ", err)
	}
	err = commonTestHelper.Pull()
	if err != nil {
		fmt.Println("Git pull code fail ", err)
	}
	err = commonTestHelper.BuildCli()
	if err != nil {
		fmt.Println("Build cli fail ", err)
	}
	err = commonTestHelper.RunServer()
	if err != nil {
		fmt.Println("Build server fail ")
	}
	commonTestHelper.TestCli()
}
