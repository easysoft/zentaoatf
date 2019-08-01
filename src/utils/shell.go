package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func ExeShell(cmdStr string) (string, error) {
	var cmd *exec.Cmd
	if IsWin() {
		cmd = exec.Command(cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}

func ExecCommand(commandName string) string {
	cmd := exec.Command("/bin/bash", "-c", commandName)

	output := make([]string, 0)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return ""
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		Printt(line)
		output = append(output, line)
	}

	cmd.Wait()

	return strings.Join(output, "")
}

func ExecCommandBlock(commandName string) string {
	cmd := exec.Command("/bin/bash", "-c", commandName)
	var out bytes.Buffer
	cmd.Stdout = &out

	_ = cmd.Run()

	return out.String()
}
