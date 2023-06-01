package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	uiTest "github.com/easysoft/zentaoatf/test/helper/zentao/ui"
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

	flagSet.BoolVar(&commConsts.Verbose, "verbose", false, "")

	testing.Init()
	flagSet.Parse(os.Args[1:])

	initTest(version)
	initZentao(runFrom, version)

	doTest(testToRun)
}

func initTest(version string) (err error) {
	commConsts.ExecFrom = commConsts.FromCmd
	serverConfig.InitLog()
	serverConfig.InitExecLog(constTestHelper.RootPath)

	commConsts.ZtfDir = constTestHelper.RootPath
	i118Utils.Init("zh-CN", commConsts.AppServer)

	fmt.Println(version)

	return
}

func initZentao(runFrom, version string) (err error) {
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
		err = commonTestHelper.BuildCli()
		if err != nil {
			fmt.Println("Build cli fail ", err)
		}
	}

	return
}

func doTest(testToRun string) (err error) {
	testPath := fmt.Sprintf(`%stest`, constTestHelper.RootPath)
	if runtime.GOOS == "windows" {
		testPath = fmt.Sprintf(`%stest`, constTestHelper.RootPath)
	}

	req := serverDomain.TestSet{
		WorkspacePath: testPath,
		Cmd:           fmt.Sprintf("go test %s -v", testToRun),
		TestTool:      commConsts.GoTest,
	}
	fmt.Println(testPath, req.Cmd)

	// exec testing
	report := execSuite(req, "restapi")
	report.ProductId = 82

	// submit result for test
	if runFrom != "jenkins" {
		config := commDomain.WorkspaceConf{
			Url:      "https://back.zcorp.cc/pms",
			Username: "chenqi",
			Password: "th2ISxOVXcoUiMLazk1b"}
		err = zentaoHelper.CommitResult(report, report.ProductId, 0, config, nil)

	}

	return
}

func execSuite(req serverDomain.TestSet, unitType string) (report commDomain.ZtfReport) {
	commConsts.AllureReportDir = filepath.Join(unitType, "allure-results")

	pth := filepath.Join(req.WorkspacePath, commConsts.AllureReportDir)
	fileUtils.RmDir(pth)
	fileUtils.MkDirIfNeeded(pth)

	startTime := time.Now()

	// run
	execHelper.RunUnitTest(nil, req.Cmd, req.WorkspacePath, nil)

	entTime := time.Now()

	// gen report
	req.ResultDir = commConsts.AllureReportDir
	req.ZipDir = req.ResultDir
	if !fileUtils.IsAbsolutePath(req.ResultDir) {
		req.ResultDir = filepath.Join(req.WorkspacePath, req.ResultDir)
	}

	if !fileUtils.IsAbsolutePath(req.ZipDir) {
		req.ZipDir = filepath.Join(req.WorkspacePath, req.ZipDir)
	}

	report = execHelper.GenUnitTestReport(req, startTime.Unix(), entTime.Unix(), nil, nil)
	fmt.Printf("执行：%v, 成功：%v，失败：%v \n", report.Total, report.Pass, report.Fail)

	return report
}
