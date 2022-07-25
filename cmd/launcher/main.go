package main

import (
	"fmt"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		cleanup()
		os.Exit(0)
	}()

	pth := ""
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		pth = "gui\\ztf.exe"
		cmd = exec.Command("start", pth)
	}

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to start ztf gui, path %s, err %s", pth, err.Error())
	}
}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
