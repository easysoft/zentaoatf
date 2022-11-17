package commonTestHelper

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	analysisHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/analysis"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"

	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
)

func TestCli() (err error) {
	cmdStr := fmt.Sprintf(`%sztf allure -allureReportDir ./test/cli/allure-results go test %stest/cli -v`, RootPath, RootPath)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmdStr = fmt.Sprintf(`%sztf.exe allure -allureReportDir .\test\cli\allure-results go test %stest\cli -v`, RootPath, RootPath)
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}
	cmd.Dir = RootPath
	fmt.Println(cmd.String())

	if cmd == nil {
		err = errors.New("cmd is nil")
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}
	defer func() {
		cmd.Process.Kill()
		execHelper.KillProcessByUUID("cli_auto_test")
	}()
	reader1 := bufio.NewReader(stdout)
	reportDir := ""
	for {
		line, err2 := reader1.ReadString('\n')
		line = strings.TrimSpace(line)
		fmt.Println(line)
		lineRune := []rune(line)
		if len(lineRune) >= 6 && string(lineRune[:6]) == "Report" {
			reportDir = strings.TrimSpace(string(lineRune[6 : len(lineRune)-1]))
			break
		} else if len(lineRune) >= 6 && string(lineRune[:2]) == "报告" {
			reportDir = strings.TrimSpace(string(lineRune[2 : len(lineRune)-1]))
			break
		}
		if err2 != nil {
			err = err2
			return
		}
		if err != nil || io.EOF == err {
			break
		}
	}
	if reportDir == "" {
		return
	}

	report, err := analysisHelper.ReadReportByPath(strings.Replace(reportDir, "result.txt", "result.json", 1))
	if err != nil {
		return
	}
	config := commDomain.WorkspaceConf{Url: "http://127.0.0.1:8081/", Password: "Test123456.", Username: "admin"}

	err = zentaoHelper.CommitResult(report, 1, 0, config, nil)
	return
}

func TestUi() (err error) {
	testPath := fmt.Sprintf(`%stest`, RootPath)
	if runtime.GOOS == "windows" {
		testPath = fmt.Sprintf(`%stest`, RootPath)
	}
	req := serverDomain.TestSet{
		WorkspacePath: testPath,
		Cmd:           "go test ./ui -v",
		TestTool:      commConsts.GoTest,
	}
	report := ExecUnit(req)

	config := commDomain.WorkspaceConf{Url: "http://127.0.0.1:8081/", Password: "Test123456.", Username: "admin"}

	err = zentaoHelper.CommitResult(report, 1, 0, config, nil)
	return
}

func ExecUnit(
	req serverDomain.TestSet) (report commDomain.ZtfReport) {
	commConsts.AllureReportDir = "ui/allure-results"
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
	return report
}
