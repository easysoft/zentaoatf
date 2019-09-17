package shellUtils

import (
	"bufio"
	"bytes"
	"fmt"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	"io"
	"os/exec"
	"path"
	"strings"
)

func ExeShell(cmdStr string) (string, error) {
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		cmd = exec.Command(cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}

func ExecFile(filePath string) string {
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		ext := path.Ext(filePath)
		if ext == ".php" {
			filePath = "php " + filePath
		}

		cmd = exec.Command("cmd", "/C", filePath)
	} else {
		exec.Command("/bin/bash", "-c", "chmod +x "+filePath)

		cmd = exec.Command("/bin/bash", "-c", "php "+filePath)
	}

	output := make([]string, 0)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return ""
	}

	cmd.Start()

	if err != nil {
		return fmt.Sprint(err)
	}

	reader := bufio.NewReader(stdout)

	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		//logUtils.Trace(strings.TrimRight(line, "\n"))
		output = append(output, line)
	}

	cmd.Wait()

	return strings.Join(output, "")
}

func ExecFile2(commandName string) string {
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", commandName)
	} else {
		commandName = "chmod +x " + commandName + "; " + commandName + ";"
		cmd = exec.Command("/bin/bash", "-c", commandName)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	var ret string
	err := cmd.Run()
	if err != nil {
		ret = fmt.Sprint(err) + " : " + stderr.String()
	} else {
		ret = out.String()
	}

	return ret
}
