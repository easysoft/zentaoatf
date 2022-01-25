package scriptUtils

import (
	"bufio"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	dateUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/date"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	serverLog "github.com/aaronchen2k/deeptest/internal/server/core/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/config"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/script"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/zentao"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"github.com/mattn/go-runewidth"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

func Exec(ch chan int, sendOutputMsg, sendExecMsg func(info, isRunning string, msg websocket.Message), req serverDomain.WsReq, msg websocket.Message) (
	err error) {

	serverLog.InitExecLog(req.ProjectPath)

	if req.Act == commConsts.ExecCase {
		ExecCase(ch, sendOutputMsg, sendExecMsg, req, msg)
	} else if req.Act == commConsts.ExecModule {
		ExecModule(ch, sendOutputMsg, sendExecMsg, req, msg)
	} else if req.Act == commConsts.ExecSuite {
		ExecSuite(ch, sendOutputMsg, sendExecMsg, req, msg)
	} else if req.Act == commConsts.ExecTask {
		ExecTask(ch, sendOutputMsg, sendExecMsg, req, msg)
	}

	return
}

func ExecCase(ch chan int, sendOutputMsg, sendExecMsg func(info, isRunning string, msg websocket.Message), req serverDomain.WsReq, msg websocket.Message) (report commDomain.ZtfReport, pathMaxWidth int, err error) {
	cases := req.Cases
	return Run(ch, sendOutputMsg, sendExecMsg, req.ProjectPath, cases, msg)
}

func ExecModule(ch chan int, sendOutputMsg, sendExecMsg func(info, isRunning string, msg websocket.Message), req serverDomain.WsReq, msg websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {
	cases := zentaoUtils.GetCasesByModule(req.ProductId, req.ModuleId, req.ProjectPath)
	return Run(ch, sendOutputMsg, sendExecMsg, req.ProjectPath, cases, msg)
}

func ExecSuite(ch chan int, sendOutputMsg, sendExecMsg func(info, isRunning string, msg websocket.Message), req serverDomain.WsReq, msg websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {
	cases := zentaoUtils.GetCasesBySuite(req.ProductId, req.SuiteId, req.ProjectPath)
	return Run(ch, sendOutputMsg, sendExecMsg, req.ProjectPath, cases, msg)
}

func ExecTask(ch chan int, sendOutputMsg, sendExecMsg func(info, isRunning string, msg websocket.Message), req serverDomain.WsReq, msg websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {
	cases := zentaoUtils.GetCasesByTask(req.ProductId, req.TaskId, req.ProjectPath)
	return Run(ch, sendOutputMsg, sendExecMsg, req.ProjectPath, cases, msg)
}

func Run(ch chan int, sendOutputMsg, sendExecMsg func(info, isRunning string, msg websocket.Message), projectPath string, cases []string, msg websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {

	conf := configUtils.LoadByProjectPath(projectPath)

	casesToRun, casesToIgnore := filterCases(cases, conf)

	numbMaxWidth := 0
	numbMaxWidth, pathMaxWidth = getNumbMaxWidth(casesToRun)
	report = genReport()

	ExeScripts(casesToRun, casesToIgnore, projectPath, conf, &report, pathMaxWidth, numbMaxWidth, ch, sendOutputMsg, sendExecMsg, msg)
	GenZTFTestReport(report, pathMaxWidth, projectPath, sendOutputMsg, sendExecMsg, msg)

	return
}

func ExeScripts(casesToRun []string, casesToIgnore []string, projectPath string, conf commDomain.ProjectConf,
	report *commDomain.ZtfReport, pathMaxWidth int, numbMaxWidth int, ch chan int,
	sendOutputMsg, sendExecMsg func(info, isRunning string, msg websocket.Message), wsMsg websocket.Message) {

	now := time.Now()
	startTime := now.Unix()
	report.StartTime = startTime

	postFix := ":"
	if len(casesToRun) == 0 {
		postFix = "."
	}

	temp := i118Utils.Sprintf("found_scripts", strconv.Itoa(len(casesToRun))) + postFix
	sendExecMsg(temp, "", wsMsg)
	logUtils.ExecConsolef(color.FgCyan, temp)
	logUtils.ExecResult(temp)

	if len(casesToIgnore) > 0 {
		temp := i118Utils.Sprintf("ignore_scripts", strconv.Itoa(len(casesToIgnore))) + postFix
		sendExecMsg(temp, "", wsMsg)
		logUtils.ExecConsolef(color.FgCyan, temp)
		logUtils.ExecResult(temp)
	}

	for idx, file := range casesToRun {
		ExeScript(file, projectPath, conf, report, idx, len(casesToRun), pathMaxWidth, numbMaxWidth, ch, sendOutputMsg, sendExecMsg, wsMsg)

		select {
		case <-ch:
			msg := i118Utils.Sprintf("exit_exec_all")
			sendExecMsg(msg, "", wsMsg)
			logUtils.ExecConsolef(color.FgCyan, msg)
			logUtils.ExecFilef(msg)

			goto ExitAllCase
		default:
		}
	}

ExitAllCase:
	endTime := time.Now().Unix()
	report.EndTime = endTime
	report.Duration = endTime - startTime
}

func ExeScript(scriptFile, projectPath string, conf commDomain.ProjectConf, report *commDomain.ZtfReport, idx,
	total, pathMaxWidth, numbMaxWidth int,
	ch chan int, sendOutputMsg, sendExecMsg func(s, isRunning string, msg websocket.Message), wsMsg websocket.Message) {

	startTime := time.Now()

	startMsg := i118Utils.Sprintf("start_execution", scriptFile, dateUtils.DateTimeStr(startTime))
	sendExecMsg(startMsg, "", wsMsg)
	logUtils.ExecConsolef(-1, startMsg)
	logUtils.ExecFilef(startMsg)

	logs := ""
	stdOutput, errOutput := RunScript(scriptFile, projectPath, conf, ch, sendOutputMsg, sendExecMsg, wsMsg)
	stdOutput = strings.Trim(stdOutput, "\n")

	if stdOutput != "" {
		logs = stdOutput
	}
	if errOutput != "" {
		sendOutputMsg(errOutput, "", wsMsg)
		logUtils.ExecConsolef(-1, errOutput)
		logUtils.ExecFilef(errOutput)
	}

	entTime := time.Now()
	secs := fmt.Sprintf("%.2f", float32(entTime.Sub(startTime)/time.Second))

	endMsg := i118Utils.Sprintf("end_execution", scriptFile, dateUtils.DateTimeStr(entTime))
	sendExecMsg(endMsg, "", wsMsg)
	logUtils.ExecConsolef(-1, endMsg)
	logUtils.ExecFilef(endMsg)

	CheckCaseResult(scriptFile, logs, report, idx, total, secs, pathMaxWidth, numbMaxWidth,
		sendOutputMsg, sendExecMsg, wsMsg)

	if idx < total-1 {
		logUtils.Infof("")
	}
}

func RunScript(filePath, projectPath string, conf commDomain.ProjectConf,
	ch chan int, sendOutputMsg, sendExecMsg func(s, isRunning string, wsMsg websocket.Message), wsMsg websocket.Message) (
	stdOutput string, errOutput string) {

	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		lang := langUtils.GetLangByFile(filePath)

		scriptInterpreter := ""
		if strings.ToLower(lang) != "bat" {
			scriptInterpreter = configUtils.GetFieldVal(conf, stringUtils.UcFirst(lang))
		}
		if scriptInterpreter != "" {
			if strings.Index(strings.ToLower(scriptInterpreter), "autoit") > -1 {
				cmd = exec.Command("cmd", "/C", scriptInterpreter, filePath, "|", "more")
			} else {
				cmd = exec.Command("cmd", "/C", scriptInterpreter, filePath)
			}
		} else if strings.ToLower(lang) == "bat" {
			cmd = exec.Command("cmd", "/C", filePath)
		} else {
			msg := i118Utils.I118Prt.Sprintf("no_interpreter_for_run", lang, filePath)
			sendOutputMsg(msg, "", wsMsg)
			logUtils.ExecConsolef(-1, msg)
			logUtils.ExecFilef(msg)
		}
	} else {
		err := os.Chmod(filePath, 0777)
		if err != nil {
			msg := i118Utils.I118Prt.Sprintf("exec_cmd_fail", filePath, err.Error())
			sendOutputMsg(msg, "", wsMsg)
			logUtils.ExecConsolef(-1, msg)
			logUtils.ExecFilef(msg)
		}

		filePath = "\"" + filePath + "\""
		cmd = exec.Command("/bin/bash", "-c", filePath)
	}

	cmd.Dir = projectPath

	if cmd == nil {
		msgStr := i118Utils.Sprintf("cmd_empty")

		sendOutputMsg(msgStr, "", wsMsg)
		logUtils.ExecConsolef(color.FgRed, msgStr)
		logUtils.ExecFilef(msgStr)

		return "", msgStr
	}

	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		sendOutputMsg(err1.Error(), "", wsMsg)
		logUtils.ExecConsolef(color.FgRed, err1.Error())
		logUtils.ExecFilef(err1.Error())

		return "", err1.Error()
	} else if err2 != nil {
		sendOutputMsg(err2.Error(), "", wsMsg)
		logUtils.ExecConsolef(color.FgRed, err2.Error())
		logUtils.ExecFilef(err2.Error())

		return "", err2.Error()
	}

	cmd.Start()

	isTerminal := false
	reader1 := bufio.NewReader(stdout)
	stdOutputArr := make([]string, 0)
	for {
		line, err2 := reader1.ReadString('\n')
		if line != "" {
			sendOutputMsg(line, "", wsMsg)
			logUtils.ExecConsole(1, line)
			logUtils.ExecFile(line)

			isTerminal = true
		}

		if err2 != nil || io.EOF == err2 {
			break
		}
		stdOutputArr = append(stdOutputArr, line)

		select {
		case <-ch:
			msg := i118Utils.Sprintf("exit_exec_curr")
			sendExecMsg(msg, "", wsMsg)
			logUtils.ExecConsolef(color.FgCyan, msg)
			logUtils.ExecFilef(msg)

			goto ExitCurrCase
		default:
		}
	}

	cmd.Wait()

ExitCurrCase:
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

	stdOutput = strings.Join(stdOutputArr, "")
	errOutput = strings.Join(errOutputArr, "")
	return
}

func filterCases(cases []string, conf commDomain.ProjectConf) (casesToRun, casesToIgnore []string) {
	for _, cs := range cases {
		if commonUtils.IsWin() {
			if path.Ext(cs) == ".sh" { // filter by os
				continue
			}

			ext := path.Ext(cs)
			if ext != "" {
				ext = ext[1:]
			}
			lang := langUtils.ScriptExtToNameMap[ext]
			interpreter := configUtils.GetFieldVal(conf, stringUtils.UcFirst(lang))
			if interpreter == "-" || interpreter == "" {
				interpreter = ""
				casesToIgnore = append(casesToIgnore, cs)
			}
			if lang != "bat" && interpreter == "" { // ignore the ones with no interpreter set
				continue
			}

		} else if !commonUtils.IsWin() { // filter by os
			if path.Ext(cs) == ".bat" {
				continue
			}
		}

		casesToRun = append(casesToRun, cs)
	}

	return
}

func genReport() (report commDomain.ZtfReport) {
	report = commDomain.ZtfReport{TestEnv: commonUtils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, FuncResult: make([]commDomain.FuncResult, 0)}
	report.TestType = consts.TestFunc
	report.TestFrame = commConsts.AppServer

	return
}

func getNumbMaxWidth(casesToRun []string) (numbMaxWidth, pathMaxWidth int) {
	for _, cs := range casesToRun {
		lent := runewidth.StringWidth(cs)
		if lent > pathMaxWidth {
			pathMaxWidth = lent
		}

		content := fileUtils.ReadFile(cs)
		caseId := scriptUtils.ReadCaseId(content)
		if len(caseId) > numbMaxWidth {
			numbMaxWidth = len(caseId)
		}
	}

	return
}
