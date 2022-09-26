package commonTestHelper

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"

	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	analysisHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/analysis"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
)

func TestCli() (err error) {
	cmdStr := fmt.Sprintf(`%sztf allure -allureReportDir ./test/cli/allure-results go test %stest/cli -v`, RootPath, RootPath)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmdStr = fmt.Sprintf(`%sztf.exe allure -allureReportDir .\test\cli\allure-results go test %stest\cli -v -run=CliVersion`, RootPath, RootPath)
		fmt.Println(RootPath, cmdStr)
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr, "-uuid", "uuuuuuuuuuu")
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
	cmd.Process.Kill()
	execHelper.KillProcessByUUID("uuuuuuuuuuu")
	report, err := analysisHelper.ReadReportByPath(strings.Replace(reportDir, "result.txt", "result.json", 1))
	if err != nil {
		return
	}
	config := commDomain.WorkspaceConf{Url: "http://127.0.0.1:8081/", Password: "Test123456.", Username: "admin"}

	fmt.Println(report, config)
	err = zentaoHelper.CommitResult(report, 1, 0, config, nil)
	fmt.Println(err)
	return
}

func TestUi() (err error) {
	cmdStr := fmt.Sprintf(`%sztf allure -allureReportDir ./test/cli/allure-results go test %stest/ui -v`, RootPath, RootPath)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmdStr = fmt.Sprintf(`%sztf.exe allure -allureReportDir .\test\cli\allure-results go test %stest\ui -v -run=CliVersion`, RootPath, RootPath)
		fmt.Println(RootPath, cmdStr)
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr, "-uuid", "uuuuuuuuuuu")
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
	cmd.Process.Kill()
	execHelper.KillProcessByUUID("uuuuuuuuuuu")
	report, err := analysisHelper.ReadReportByPath(strings.Replace(reportDir, "result.txt", "result.json", 1))
	if err != nil {
		return
	}
	config := commDomain.WorkspaceConf{Url: "http://127.0.0.1:8081/", Password: "Test123456.", Username: "admin"}

	fmt.Println(report, config)
	err = zentaoHelper.CommitResult(report, 1, 0, config, nil)
	fmt.Println(err)
	return
}
