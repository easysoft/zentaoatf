package scriptUtils

import (
	"bufio"
	"errors"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	dateUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/date"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/server/modules/helper/zentao"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"io"
	"os/exec"
	"strings"
	"time"
)

func ExecUnit(ch chan int, sendOutputMsg,
	sendExecMsg func(info, isRunning string, wsMsg websocket.Message),
	req serverDomain.WsReq, wsMsg websocket.Message) (resultDir string, err error) {

	startTime := time.Now()
	startMsg := i118Utils.Sprintf("start_execution", req.Cmd, dateUtils.DateTimeStr(startTime))

	if commConsts.ComeFrom != "cmd" {
		sendExecMsg(startMsg, "", wsMsg)
	}

	logUtils.ExecConsolef(-1, startMsg)
	logUtils.ExecFilef(startMsg)

	RunUnitTest(ch, sendOutputMsg, sendExecMsg, req.Cmd, req.ProjectPath, wsMsg)

	entTime := time.Now()
	endMsg := i118Utils.Sprintf("end_execution", req.Cmd, dateUtils.DateTimeStr(entTime))

	if commConsts.ComeFrom != "cmd" {
		sendExecMsg(endMsg, "false", wsMsg)
	}

	logUtils.ExecConsolef(-1, endMsg)
	logUtils.ExecFilef(endMsg)

	report := GenUnitTestReport(req, startTime.Unix(), entTime.Unix(), ch, sendOutputMsg, sendExecMsg, wsMsg)
	logUtils.Infof("#v", report)

	if req.SubmitResult {
		err = zentaoUtils.CommitResult(report, req.ProductId, "0", req.ProjectPath)
	}

	return
}

func RunUnitTest(ch chan int, sendOutputMsg, sendExecMsg func(info, isRunning string, msg websocket.Message),
	cmdStr, projectPath string, wsMsg websocket.Message) (err error) {

	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	cmd.Dir = projectPath

	if cmd == nil {
		msgStr := i118Utils.Sprintf("cmd_empty")

		sendOutputMsg(msgStr, "", wsMsg)
		logUtils.ExecConsolef(color.FgRed, msgStr)
		logUtils.ExecFilef(msgStr)

		err = errors.New(msgStr)
		return
	}

	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		sendOutputMsg(err1.Error(), "", wsMsg)
		logUtils.ExecConsolef(color.FgRed, err1.Error())
		logUtils.ExecFilef(err1.Error())

		return
	} else if err2 != nil {
		sendOutputMsg(err2.Error(), "", wsMsg)
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
			sendOutputMsg(line, "", wsMsg)
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

			if commConsts.ComeFrom != "cmd" {
				sendExecMsg(msg, "", wsMsg)
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
		sendOutputMsg(errOutput, "", wsMsg)
		logUtils.ExecConsolef(-1, errOutput)
		logUtils.ExecFilef(errOutput)
	}

	return
}
