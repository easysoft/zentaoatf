package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/ergoapi/util/zos"
	"github.com/fatih/color"
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
	if !zos.IsUnix() {
		pth = "gui\\ztf.exe"
		cmd = exec.Command("cmd", "/C", "start", pth)
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
