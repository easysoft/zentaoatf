//go:build windows
// +build windows

package shellUtils

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

type window struct {
	Row uint16
	Col uint16
}

func WindowSize() window {
	win := window{0, 0}

	cmd1 := exec.Command("cmd")
	cmd1.SysProcAttr = &syscall.SysProcAttr{CmdLine: "/c mode con", HideWindow: true}

	out, _ := cmd1.Output()
	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		if win.Row > 0 && win.Col > 0 {
			return win
		}
		line = strings.TrimSpace(line)
		if strings.Contains(line, "行") || strings.Contains(line, "Row") {
			re := regexp.MustCompile(`\d+`)
			rs := re.FindAllString(line, -1)
			win.Row, _ = strconv.Atoi(rs[1])
		}
		if strings.Contains(line, "列") || strings.Contains(line, "Col") {
			re := regexp.MustCompile(`\d+`)
			rs := re.FindAllString(line, -1)
			win.Col, _ = strconv.Atoi(rs[1])
		}
	}

	return win
}
