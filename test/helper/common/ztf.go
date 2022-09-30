package commonTestHelper

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func BuildCli() (err error) {
	outPath := fmt.Sprintf("%s%s", RootPath, "ztf")
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
	cmd.Dir = RootPath
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
	cmd.Dir = RootPath
	fmt.Println(cmd.String())
	err = cmd.Start()
	if err != nil {
		return
	}
	return
}

func RunUi() (err error) {
	var cmd *exec.Cmd
	cmd = exec.Command("npm", "run", "serve", "-uuid=ui_auto_test")
	cmd.Dir = RootPath + FilePthSep + "ui"
	fmt.Println(cmd.String())
	err = cmd.Start()
	if err != nil {
		return
	}
	return
}

func GetZtfPath() string {
	ztfPath := fmt.Sprintf("%s%s", RootPath, "ztf")
	if runtime.GOOS == "windows" {
		ztfPath += ".exe"
	}
	return ztfPath
}

func GetZtfProductPath() string {
	return fmt.Sprintf("%s%s%s%s%s%s%s%s", RootPath, "test", FilePthSep, "demo", FilePthSep, "php", FilePthSep, "product1")
}

func GetPhpWorkspacePath() string {
	return fmt.Sprintf("%s%s%s%s%s%s%s", RootPath, "test", FilePthSep, "demo", FilePthSep, "php", FilePthSep)
}
