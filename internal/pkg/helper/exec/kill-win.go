//go:build windows
// +build windows

package execHelper

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func KillProcessByUUID(uuid string) {
	cmd1 := exec.Command("cmd")
	cmd1.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c WMIC path win32_process where "CommandLine like '%%%s%%'" get Processid,Caption`, uuid), HideWindow: true}

	out, _ := cmd1.Output()
	lines := strings.Split(string(out), "\n")
	for index, line := range lines {
		if index == 0 {
			continue
		}
		line = strings.TrimSpace(line)
		cols := strings.Split(line, " ")
		pid := ""
		if len(cols) > 3 {
			pid = cols[3]
		} else if len(cols) > 2 {
			pid = cols[2]
		}
		if pid != "" {
			fmt.Println(fmt.Sprintf(`taskkill /F /pid %s`, pid))
			cmd2 := exec.Command("cmd")
			cmd2.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c taskkill /F /pid %s`, pid), HideWindow: true}
			cmd2.Start()
		}
	}
}
