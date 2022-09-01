//go:build !windows
// +build !windows

package execHelper

import (
	"fmt"
	"os/exec"
)

func KillProcessByUUID(uuid string) {
	command := fmt.Sprintf(`ps -ef | grep %s | grep -v "grep" | awk '{print $2}' | xargs kill -9`, uuid)
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Start()
}
