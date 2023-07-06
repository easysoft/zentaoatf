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
	shellUtils "github.com/easysoft/zentaoatf/pkg/lib/shell"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	dateUtils "github.com/easysoft/zentaoatf/pkg/lib/date"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

func ExecUnit(ch chan int, req serverDomain.TestSet, wsMsg *websocket.Message) (err error) {
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
		if fileUtils.GetExtName(req.ResultDir) == "" { // not file
			fileUtils.MkDirIfNeeded(req.ResultDir)
		}

		RunUnitTest(ch, req.Cmd, req.WorkspacePath, wsMsg)

		// end msg
		endMsg := i118Utils.Sprintf("end_execution", req.Cmd, dateUtils.DateTimeStr(entTime))
		logUtils.ExecConsolef(-1, endMsg+"\n")
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
		report.JacocoResult = GenJacocoCovReport()
	}

	// submit result
	if req.SubmitResult && (report.FuncResult != nil || report.UnitResult != nil) {
		workspaceDir := req.WorkspacePath
		if commConsts.ExecFrom == commConsts.FromCmd {
			workspaceDir = commConsts.ZtfDir
		}

		config := configHelper.LoadByWorkspacePath(workspaceDir)

		err = zentaoHelper.CommitResult(report, req.ProductId, req.TaskId, config, wsMsg)
	}

	if commConsts.ExecFrom == commConsts.FromClient {
		websocketHelper.SendExecMsg("", "false", commConsts.Run, nil, wsMsg)
	}

	return
}

func RunUnitTest(ch chan int, cmdStr, workspacePath string, wsMsg *websocket.Message) (err error) {
	key := stringUtils.Md5(workspacePath)

	cmd := shellUtils.GetCmd(cmdStr)

	if cmd == nil {
		msgStr := i118Utils.Sprintf("cmd_empty")
		websocketHelper.SendOutputMsgIfNeed(msgStr, "", iris.Map{"key": key}, wsMsg)

		logUtils.ExecConsolef(color.FgRed, msgStr)
		logUtils.ExecFilef(msgStr)

		err = errors.New(msgStr)
		return
	}

	cmd.Dir = workspacePath
	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		websocketHelper.SendOutputMsgIfNeed(err1.Error(), "", iris.Map{"key": key}, wsMsg)
		logUtils.ExecConsolef(color.FgRed, err1.Error())
		logUtils.ExecFilef(err1.Error())

		return
	} else if err2 != nil {
		websocketHelper.SendOutputMsgIfNeed(err2.Error(), "", iris.Map{"key": key}, wsMsg)
		logUtils.ExecConsolef(color.FgRed, err2.Error())
		logUtils.ExecFilef(err2.Error())

		return
	}

	cmd.Start()

	isTerminal := printStdout(stdout, ch, cmd, key, wsMsg)

	printStderr(isTerminal, stderr, key, wsMsg)

	cmd.Wait()

	return
}

func printStdout(stdout io.ReadCloser, ch chan int, cmd *exec.Cmd, key string, wsMsg *websocket.Message) (isTerminal bool) {
	isTerminal = false
	reader1 := bufio.NewReader(stdout)

	for {
		line, err3 := reader1.ReadString('\n')
		if line != "" {
			line = stringUtils.Convert2Utf8IfNeeded(line)
			websocketHelper.SendOutputMsgIfNeed(line, "", iris.Map{"key": key}, wsMsg)
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

			websocketHelper.SendExecMsgIfNeed(msg, "", commConsts.Run, nil, wsMsg)

			logUtils.ExecConsolef(color.FgCyan, msg)
			logUtils.ExecFilef(msg)

			return
		default:
		}
	}

	return
}

func printStderr(isTerminal bool, stderr io.ReadCloser, key string, wsMsg *websocket.Message) {
	errOutputArr := make([]string, 0)
	if !isTerminal {
		reader2 := bufio.NewReader(stderr)

		for {
			line, err2 := reader2.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}
			line = stringUtils.Convert2Utf8IfNeeded(line)
			errOutputArr = append(errOutputArr, line)
		}
	}

	errOutput := strings.Join(errOutputArr, "")

	if errOutput != "" {
		websocketHelper.SendOutputMsgIfNeed(errOutput, "", iris.Map{"key": key}, wsMsg)
		logUtils.ExecConsolef(-1, errOutput)
		logUtils.ExecFilef(errOutput)
	}
}
