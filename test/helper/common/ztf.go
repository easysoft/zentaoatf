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

func BuildServer() (err error) {
	outPath := fmt.Sprintf("%s%s", RootPath, "server")
	cliPath := `./cmd/server/main.go`
	if runtime.GOOS == "windows" {
		cliPath = `.\cmd\server\main.go`
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

func GetZtfPath() string {
	ztfPath := fmt.Sprintf("%s%s", RootPath, "ztf")
	if runtime.GOOS == "windows" {
		ztfPath += ".exe"
	}
	return ztfPath
}
