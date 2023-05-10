//go:build windows
// +build windows

package execHelper

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

type window struct {
	Row int
	Col int
}

func WindowSize(uuid string) window {
	win := window{0, 0}

	cmd1 := exec.Command("cmd")
	cmd1.SysProcAttr = &syscall.SysProcAttr{CmdLine: "/c mode con", HideWindow: true}

	out, _ := cmd1.Output()
	lines := strings.Split(string(out), "\n")

	for index, line := range lines {
		if win.Row > 0 && win.Col > 0 {
			return win
		}
		line = strings.TrimSpace(line)
		if strings.Contain("行") || strings.Contain("Row") {
			re := regexp.MustCompile(`\d+`)
			rs := re.FindAllString(out.String(), -1)
			win.Row, _ = strconv.Atoi(rs[1])
		}
		if strings.Contain("列") || strings.Contain("Col") {
			re := regexp.MustCompile(`\d+`)
			rs := re.FindAllString(out.String(), -1)
			win.Col, _ = strconv.Atoi(rs[1])
		}
	}

	return win
}
