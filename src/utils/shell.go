package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func ExeShell(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}

func ExecCommand(commandName string) []string {
	cmd := exec.Command("/bin/bash", "-c", commandName)

	output := make([]string, 0)

	//显示运行的命令
	fmt.Println("begin to run " + strings.Join(cmd.Args, " "))

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return output
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
		output = append(output, line)
	}

	cmd.Wait()
	return output
}
