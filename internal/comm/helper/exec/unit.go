package execHelper

import (
	"bufio"
	"errors"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	websocketHelper "github.com/easysoft/zentaoatf/internal/comm/helper/websocket"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/comm/helper/zentao"
	commonUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/common"
	dateUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/date"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"io"
	"os/exec"
	"strings"
	"time"
)

func ExecUnit(ch chan int,
	req serverDomain.TestSet, wsMsg *websocket.Message) (resultDir string, err error) {

	startTime := time.Now()
	startMsg := i118Utils.Sprintf("start_execution", req.Cmd, dateUtils.DateTimeStr(startTime))

	if commConsts.ExecFrom != commConsts.FromCmd {
		websocketHelper.SendExecMsg(startMsg, "", commConsts.Run, wsMsg)
	}

	logUtils.ExecConsolef(-1, startMsg)
	logUtils.ExecFilef(startMsg)

	RunUnitTest(ch, req.Cmd, req.WorkspacePath, wsMsg)

	entTime := time.Now()
	endMsg := i118Utils.Sprintf("end_execution", req.Cmd, dateUtils.DateTimeStr(entTime))

	if commConsts.ExecFrom != commConsts.FromCmd {
		websocketHelper.SendExecMsg(endMsg, "false", commConsts.Run, wsMsg)
	}

	logUtils.ExecConsolef(-1, endMsg)
	logUtils.ExecFilef(endMsg)

	report := GenUnitTestReport(req, startTime.Unix(), entTime.Unix(), ch, wsMsg)

	if req.SubmitResult {
		config := configHelper.LoadByWorkspacePath(req.WorkspacePath)
		err = zentaoHelper.CommitResult(report, req.ProductId, 0, config)
	}

	return
}

func RunUnitTest(ch chan int, cmdStr, workspacePath string, wsMsg *websocket.Message) (err error) {
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	cmd.Dir = workspacePath

	if cmd == nil {
		msgStr := i118Utils.Sprintf("cmd_empty")
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendOutputMsg(msgStr, "", wsMsg)
		}
		logUtils.ExecConsolef(color.FgRed, msgStr)
		logUtils.ExecFilef(msgStr)

		err = errors.New(msgStr)
		return
	}

	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendOutputMsg(err1.Error(), "", wsMsg)
		}
		logUtils.ExecConsolef(color.FgRed, err1.Error())
		logUtils.ExecFilef(err1.Error())

		return
	} else if err2 != nil {
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendOutputMsg(err2.Error(), "", wsMsg)
		}
		logUtils.ExecConsolef(color.FgRed, err2.Error())
		logUtils.ExecFilef(err2.Error())

		return
	}

	cmd.Start()

	isTerminal := false
	reader1 := bufio.NewReader(stdout)
	for {
		line, err3 := reader1.ReadString('\n')
		if line != "" {
			if commConsts.ExecFrom != commConsts.FromCmd {
				websocketHelper.SendOutputMsg(line, "", wsMsg)
			}
			logUtils.ExecConsole(1, line)
			logUtils.ExecFile(line)

			isTerminal = true
		}

		if err3 != nil || io.EOF == err3 {
			break
		}

		select {
		case <-ch:
			msg := i118Utils.Sprintf("exit_exec_curr")

			if commConsts.ExecFrom != commConsts.FromCmd {
				websocketHelper.SendExecMsg(msg, "", commConsts.Run, wsMsg)
			}

			logUtils.ExecConsolef(color.FgCyan, msg)
			logUtils.ExecFilef(msg)

			goto ExitUnitTest
		default:
		}
	}

	cmd.Wait()

ExitUnitTest:
	errOutputArr := make([]string, 0)
	if !isTerminal {
		reader2 := bufio.NewReader(stderr)

		for {
			line, err2 := reader2.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}
			errOutputArr = append(errOutputArr, line)
		}
	}

	errOutput := strings.Join(errOutputArr, "")

	if errOutput != "" {
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendOutputMsg(errOutput, "", wsMsg)
		}
		logUtils.ExecConsolef(-1, errOutput)
		logUtils.ExecFilef(errOutput)
	}

	return
}
