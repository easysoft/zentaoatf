package commonTestHelper

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	uiTest "github.com/easysoft/zentaoatf/cmd/test/helper/zentao/ui"
)

func BuildCli() (err error) {
	outPath := fmt.Sprintf("%s%s", constTestHelper.RootPath, "ztf")
	cliPath := `./cmd/command/main.go`
	if runtime.GOOS == "windows" {
		cliPath = `.\cmd\command\main.go`
		outPath += ".exe"
	}

	_, err = os.Stat(outPath)
	if err != nil && os.IsExist(err) {
		os.Remove(outPath)
	}

	var cmd *exec.Cmd
	cmd = exec.Command("go", "build", "-o", outPath, cliPath)
	cmd.Dir = constTestHelper.RootPath
	fmt.Println(cmd.String())

	_, err = cmd.CombinedOutput()
	if err != nil {
		return
	}

	return
}

func RunServer() (err error) {
	ztfPath := GetZtfPath()

	var cmd *exec.Cmd
	cmd = exec.Command(ztfPath, "-P", "8085", "-uuid=ui_auto_test")
	cmd.Dir = constTestHelper.RootPath

	fmt.Println(cmd.String())
	err = cmd.Start()
	if err != nil {
		return
	}

	return
}

func RunUi() (err error) {
	var cmd *exec.Cmd
	cmd = exec.Command("yarn", "run", "serve", "--port", strconv.Itoa(constTestHelper.UiPort), "-uuid=ui_auto_test")
	cmd.Dir = constTestHelper.RootPath + constTestHelper.FilePthSep + "ui"

	fmt.Println(cmd.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	reader1 := bufio.NewReader(stdout)
	go func() {
		for {
			line, err2 := reader1.ReadString('\n')
			line = strings.TrimSpace(line)
			fmt.Println(line)
			if err2 != nil {
				return
			}

			if err != nil || io.EOF == err {
				break
			}
		}
	}()

	WaitZtfAccessed()
	return
}

func GetZtfPath() string {
	ztfPath := fmt.Sprintf("%s%s", constTestHelper.RootPath, "ztf")
	if runtime.GOOS == "windows" {
		ztfPath += ".exe"
	}
	return ztfPath
}

func GetZtfProductPath() string {
	return filepath.Join(constTestHelper.RootPath, "cmd", "test", "demo", "php", "product1")
}

func GetPhpWorkspacePath() string {
	return filepath.Join(constTestHelper.RootPath, "cmd", "test", "demo", "php")
}

func WaitZtfAccessed() {
	isTimeout := false
	time.AfterFunc(120*time.Second, func() {
		isTimeout = true
	})

	for {
		status := uiTest.GetStatus(constTestHelper.ZtfUrl)
		if isTimeout || status {
			return
		}

		time.Sleep(time.Second)
	}
}
