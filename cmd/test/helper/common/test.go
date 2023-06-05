package commonTestHelper

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"

	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
)

func TestCli() (err error) {
	testPath := fmt.Sprintf(`%stest`, constTestHelper.RootPath)
	if runtime.GOOS == "windows" {
		testPath = fmt.Sprintf(`%stest`, constTestHelper.RootPath)
	}
	req := serverDomain.TestSet{
		WorkspacePath: testPath,
		Cmd:           "go test ./cli/cli_set_test.go && go test ./cli -v",
		TestTool:      commConsts.GoTest,
	}
	fmt.Println(testPath, req.Cmd)
	report := ExecUnit(req, "cli")

	config := commDomain.WorkspaceConf{Url: constTestHelper.ZentaoSiteUrl + "/", Password: constTestHelper.ZentaoPassword, Username: constTestHelper.ZentaoUsername}

	err = zentaoHelper.CommitResult(report, 1, 0, config, nil)

	if report.Fail > 0 {
		os.Exit(1)
	}

	return
}

func TestUi() (err error) {
	var screenshotPath = fmt.Sprintf("%stest/screenshot", constTestHelper.RootPath)
	os.RemoveAll(screenshotPath)
	fileUtils.MkDirIfNeeded(screenshotPath)
	testPath := fmt.Sprintf(`%stest`, constTestHelper.RootPath)
	if runtime.GOOS == "windows" {
		testPath = fmt.Sprintf(`%stest`, constTestHelper.RootPath)
	}
	req := serverDomain.TestSet{
		WorkspacePath: testPath,
		Cmd:           "go test ./ui -v -timeout 10m",
		TestTool:      commConsts.GoTest,
	}
	report := ExecUnit(req, "ui")

	config := commDomain.WorkspaceConf{Url: constTestHelper.ZentaoSiteUrl + "/", Password: constTestHelper.ZentaoPassword, Username: constTestHelper.ZentaoUsername}

	err = zentaoHelper.CommitResult(report, 1, 0, config, nil)

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
