package commonTestHelper

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"

	dateUtils "github.com/easysoft/zentaoatf/pkg/lib/date"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
)

func TestCli(version string) (err error) {
	testPath := fmt.Sprintf(`%scmd/test`, constTestHelper.RootPath)
	if runtime.GOOS == "windows" {
		testPath = fmt.Sprintf(`%scmd\test`, constTestHelper.RootPath)
	}
	req := serverDomain.TestSet{
		WorkspacePath: testPath,
		Cmd:           "go test ./cli/cli_set_test.go && go test ./cli -v",
		TestTool:      commConsts.GoTest,
	}
	fmt.Println(testPath, req.Cmd)
	report := ExecUnit(req, "cli")
	report.ProductId = 82
	report.Name = fmt.Sprintf("禅道版本:%s Cli测试-%s", version, dateUtils.TimeStr(time.Now()))

	config := commDomain.WorkspaceConf{
		Url:      "https://back.zcorp.cc/pms",
		Username: "chenqi",
		Password: "th2ISxOVXcoUiMLazk1b"}

	err = zentaoHelper.CommitResult(report, report.ProductId, 0, config, nil)

	if report.Fail > 0 {
		os.Exit(1)
	}

	return
}

func TestUi(version string) (err error) {
	var screenshotPath = fmt.Sprintf("%scmd/test/screenshot", constTestHelper.RootPath)
	os.RemoveAll(screenshotPath)
	fileUtils.MkDirIfNeeded(screenshotPath)
	testPath := filepath.Join(constTestHelper.RootPath, "cmd", "test")

	req := serverDomain.TestSet{
		WorkspacePath: testPath,
		Cmd:           "go test ./ui/ -v -timeout 10m",
		TestTool:      commConsts.GoTest,
	}
	report := ExecUnit(req, "ui")
	report.ProductId = 82
	report.Name = fmt.Sprintf("禅道版本:%s UI测试-%s", version, dateUtils.TimeStr(time.Now()))

	config := commDomain.WorkspaceConf{
		Url:      "https://back.zcorp.cc/pms",
		Username: "chenqi",
		Password: "th2ISxOVXcoUiMLazk1b"}

	err = zentaoHelper.CommitResult(report, report.ProductId, 0, config, nil)

	if report.Fail > 0 {
		os.Exit(1)
	}
	return
}

func ExecUnit(req serverDomain.TestSet, unitType string) (report commDomain.ZtfReport) {
	if unitType == "ui" {
		commConsts.AllureReportDir = "ui/allure-results"
	} else {
		commConsts.AllureReportDir = "cli/allure-results"
	}
	pth := filepath.Join(req.WorkspacePath, commConsts.AllureReportDir)
	fileUtils.RmDir(pth)

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
	fmt.Printf("执行：%v, 成功：%v，失败：%v", report.Total, report.Pass, report.Fail)
	return report
}
