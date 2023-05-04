package execHelper

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	dateUtils "github.com/easysoft/zentaoatf/pkg/lib/date"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

func ExecUnit(ch chan int,
	req serverDomain.TestSet, wsMsg *websocket.Message) (err error) {

	key := stringUtils.Md5(req.WorkspacePath)

	// start msg
	startTime := time.Now()
	startMsg := i118Utils.Sprintf("start_execution", req.Cmd, dateUtils.DateTimeStr(startTime))
	logUtils.ExecConsolef(-1, startMsg)
	logUtils.ExecFilef(startMsg)
	if commConsts.ExecFrom == commConsts.FromClient {
		websocketHelper.SendExecMsg(startMsg, "", commConsts.Run,
			iris.Map{"key": key, "status": "start"}, wsMsg)
	}

	//deal with -allureReportDir param
	arr := strings.Split(req.Cmd, " ")
	if len(arr) > 1 && strings.TrimSpace(arr[0]) == "-allureReportDir" {
		commConsts.AllureReportDir = arr[1]
		pth := filepath.Join(req.WorkspacePath, commConsts.AllureReportDir)
		fileUtils.RmDir(pth)
		req.Cmd = strings.Join(arr[2:], " ")
	}

	getResultDirForDifferentTool(&req)

	// run
	entTime := time.Now()
	if strings.TrimSpace(req.Cmd) != "" {
		// remove old results
		fileUtils.RmDir(req.ResultDir)
		fileUtils.MkDirIfNeeded(req.ResultDir)

		RunUnitTest(ch, req.Cmd, req.WorkspacePath, wsMsg)

		// end msg
		endMsg := i118Utils.Sprintf("end_execution", req.Cmd, dateUtils.DateTimeStr(entTime))
		logUtils.ExecConsolef(-1, endMsg)
		logUtils.ExecFilef(endMsg)
		if commConsts.ExecFrom == commConsts.FromClient {
			websocketHelper.SendExecMsg(endMsg, "", commConsts.Run,
				iris.Map{"key": key, "status": "end"}, wsMsg)
		}
	}

	// gen report
	report := GenUnitTestReport(req, startTime.Unix(), entTime.Unix(), ch, wsMsg)

	// dealwith jacoco report
	if commConsts.JacocoReport != "" {
		GenJacocoCovReport()
	}

	// submit result
	if req.SubmitResult && (report.FuncResult != nil || report.UnitResult != nil) {
		configDir := req.WorkspacePath
		if commConsts.ExecFrom == commConsts.FromCmd {
			configDir = commConsts.ZtfDir
		}

		config := configHelper.LoadByWorkspacePath(configDir)

		err = zentaoHelper.CommitResult(report, req.ProductId, 0, 0, config, wsMsg)
	}

	if commConsts.ExecFrom == commConsts.FromClient {
		websocketHelper.SendExecMsg("", "false", commConsts.Run, nil, wsMsg)
	}

	return
}

func RunUnitTest(ch chan int, cmdStr, workspacePath string, wsMsg *websocket.Message) (err error) {
	key := stringUtils.Md5(workspacePath)

	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	cmd.Dir = workspacePath

	if cmd == nil {
		msgStr := i118Utils.Sprintf("cmd_empty")
		if commConsts.ExecFrom == commConsts.FromClient {
			websocketHelper.SendOutputMsg(msgStr, "", iris.Map{"key": key}, wsMsg)
		}

		logUtils.ExecConsolef(color.FgRed, msgStr)
		logUtils.ExecFilef(msgStr)

		err = errors.New(msgStr)
		return
	}

	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		if commConsts.ExecFrom == commConsts.FromClient {
			websocketHelper.SendOutputMsg(err1.Error(), "", iris.Map{"key": key}, wsMsg)
		}
		logUtils.ExecConsolef(color.FgRed, err1.Error())
		logUtils.ExecFilef(err1.Error())

		return
	} else if err2 != nil {
		if commConsts.ExecFrom == commConsts.FromClient {
			websocketHelper.SendOutputMsg(err2.Error(), "", iris.Map{"key": key}, wsMsg)
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
			if commConsts.ExecFrom == commConsts.FromClient {
				websocketHelper.SendOutputMsg(line, "", iris.Map{"key": key}, wsMsg)
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
			cmd.Process.Kill()
			msg := i118Utils.Sprintf("exit_exec_curr")

			if commConsts.ExecFrom == commConsts.FromClient {
				websocketHelper.SendExecMsg(msg, "", commConsts.Run, nil, wsMsg)
			}

			logUtils.ExecConsolef(color.FgCyan, msg)
			logUtils.ExecFilef(msg)

			goto ExitUnitTest
		default:
		}
	}

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
		if commConsts.ExecFrom == commConsts.FromClient {
			websocketHelper.SendOutputMsg(errOutput, "", iris.Map{"key": key}, wsMsg)
		}
		logUtils.ExecConsolef(-1, errOutput)
		logUtils.ExecFilef(errOutput)
	}

	cmd.Wait()

	return
}
