package shellUtils

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"github.com/ergoapi/util/zos"
)

func ExeSysCmd(cmdStr string) (string, error) {
	cmd := GetCmd(cmdStr)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	output := out.String()

	return output, err
}

func ExeShell(cmdStr string) (string, error) {
	return ExeShellInDir(cmdStr, "")
}

func ExeShellInDir(cmdStr string, dir string) (ret string, err error) {
	ret, err, _ = ExeShellInDirWithPid(cmdStr, dir)
	return
}

func ExeShellWithPid(cmdStr string) (string, error, int) {
	return ExeShellInDirWithPid(cmdStr, "")
}

func ExeShellInDirWithPid(cmdStr string, dir string) (ret string, err error, pid int) {
	cmd := GetCmd(cmdStr)
	if dir != "" {
		cmd.Dir = dir
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		logUtils.Errorf(i118Utils.Sprintf("fail_to_exec_command", cmdStr, cmd.Dir, err))
	}

	pid = cmd.Process.Pid
	ret = stringUtils.TrimAll(out.String())
	return
}

func ExeShellWithOutput(cmdStr string) ([]string, error) {
	return ExeShellWithOutputInDir(cmdStr, "")
}

func ExeShellWithOutputInDir(cmdStr string, dir string) ([]string, error) {
	return ExeShellWithEnvVarsAndOutputInDir(cmdStr, dir, nil)
}

func ExeShellWithEnvVarsAndOutputInDir(cmdStr, dir string, envVars []string) ([]string, error) {
	cmd := GetCmd(cmdStr)

	if dir != "" {
		cmd.Dir = dir
	}
	if envVars != nil && len(envVars) > 0 {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envVars...)
	}

	logUtils.InfofIfVerbose("exec cmd: %s", cmd.String())

	output := make([]string, 0)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		logUtils.InfofIfVerbose("exec cmd failed, cmd: %s, error: %s", cmd.String(), err.Error())
		return output, err
	}

	cmd.Start()

	if err != nil {
		return output, err
	}

	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		logUtils.Info(strings.TrimRight(line, "\n"))
		output = append(output, line)
	}

	cmd.Wait()

	return output, nil
}

func GetCmd(cmdStr string) (cmd *exec.Cmd) {
	if zos.IsUnix() {
		cmd = getLinuxCmd(cmdStr)
	} else {
		cmd = getWinCmd(cmdStr)
	}

	return
}
func getWinCmd(cmdStr string) (cmd *exec.Cmd) {
	return exec.Command("cmd", "/C", cmdStr)
}
func getLinuxCmd(cmdStr string) (cmd *exec.Cmd) {
	return exec.Command("/bin/bash", "-c", cmdStr)
}
