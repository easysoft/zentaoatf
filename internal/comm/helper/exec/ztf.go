package execHelper

import (
	"bufio"
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	analysisHelper "github.com/easysoft/zentaoatf/internal/comm/helper/analysis"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	langHelper "github.com/easysoft/zentaoatf/internal/comm/helper/lang"
	scriptHelper "github.com/easysoft/zentaoatf/internal/comm/helper/script"
	websocketHelper "github.com/easysoft/zentaoatf/internal/comm/helper/websocket"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/comm/helper/zentao"
	commonUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/common"
	dateUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/date"
	fileUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/string"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
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

func ExecCases(ch chan int, testSet serverDomain.TestSet, msg *websocket.Message) (report commDomain.ZtfReport, pathMaxWidth int, err error) {
	cases := testSet.Cases
	return RunZtf(ch, testSet.WorkspacePath, 0, 0, commConsts.Case, cases, msg)
}

func ExecModule(ch chan int, testSet serverDomain.TestSet, msg *websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {

	cases, err := zentaoHelper.GetCasesByModuleInDir(uint(testSet.ProductId), uint(testSet.ModuleId),
		testSet.WorkspacePath, testSet.ScriptDirParamFromCmdLine)
	if err != nil {
		return
	}

	if testSet.Seq != "" {
		cases = analysisHelper.FilterCaseByResult(cases, testSet)
	}

	return RunZtf(ch, testSet.WorkspacePath, testSet.ProductId, testSet.ModuleId, commConsts.Module, cases, msg)
}

func ExecSuite(ch chan int, testSet serverDomain.TestSet, msg *websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {
	cases, err := zentaoHelper.GetCasesBySuiteInDir(testSet.ProductId, testSet.SuiteId,
		testSet.WorkspacePath, testSet.ScriptDirParamFromCmdLine)

	if testSet.Seq != "" {
		cases = analysisHelper.FilterCaseByResult(cases, testSet)
	}

	return RunZtf(ch, testSet.WorkspacePath,
		testSet.ProductId, testSet.SuiteId, commConsts.Suite, cases, msg)
}

func ExecTask(ch chan int, testSet serverDomain.TestSet, msg *websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {
	cases, err := zentaoHelper.GetCasesByTaskInDir(testSet.ProductId, testSet.TaskId,
		testSet.WorkspacePath, testSet.ScriptDirParamFromCmdLine)
	if err != nil {
		return
	}

	if testSet.Seq != "" {
		cases = analysisHelper.FilterCaseByResult(cases, testSet)
	}

	return RunZtf(ch, testSet.WorkspacePath,
		testSet.ProductId, testSet.TaskId, commConsts.Task, cases, msg)
}

func RunZtf(ch chan int,
	workspacePath string, productId, id int, by commConsts.ExecBy, cases []string, wsMsg *websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {

	conf := configHelper.LoadByWorkspacePath(workspacePath)

	casesToRun, casesToIgnore := FilterCases(cases, conf)

	numbMaxWidth := 0
	numbMaxWidth, pathMaxWidth = getNumbMaxWidth(casesToRun)
	report = genReport(productId, id, by)

	// run
	ExeScripts(casesToRun, casesToIgnore, workspacePath, conf, &report, pathMaxWidth, numbMaxWidth, ch, wsMsg)

	// gen report
	GenZTFTestReport(report, pathMaxWidth, workspacePath, wsMsg)

	if commConsts.ExecFrom != commConsts.FromCmd {
		websocketHelper.SendExecMsg("", "false", commConsts.Run, wsMsg)
	}

	return
}

func ExeScripts(casesToRun []string, casesToIgnore []string, workspacePath string, conf commDomain.WorkspaceConf,
	report *commDomain.ZtfReport, pathMaxWidth int, numbMaxWidth int,
	ch chan int, wsMsg *websocket.Message) {

	now := time.Now()
	startTime := now.Unix()
	report.StartTime = startTime

	msg := i118Utils.Sprintf("found_scripts", len(casesToRun), workspacePath)
	if commConsts.ExecFrom != commConsts.FromCmd {
		websocketHelper.SendExecMsg(msg, "", commConsts.Run, wsMsg)
	}
	logUtils.ExecConsolef(color.FgCyan, msg)
	logUtils.ExecResult(msg)

	if len(casesToIgnore) > 0 {
		temp := i118Utils.Sprintf("ignore_scripts", strconv.Itoa(len(casesToIgnore)))
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendExecMsg(temp, "", commConsts.Run, wsMsg)
		}

		logUtils.ExecConsolef(color.FgCyan, temp)
		logUtils.ExecResult(temp)
	}

	for idx, file := range casesToRun {
		if fileUtils.IsDir(file) {
			continue
		}

		ExeScript(file, workspacePath, conf, report, idx, len(casesToRun), pathMaxWidth, numbMaxWidth, ch, wsMsg)

		select {
		case <-ch:
			msg := i118Utils.Sprintf("exit_exec_all")
			if commConsts.ExecFrom != commConsts.FromCmd {
				websocketHelper.SendExecMsg(msg, "", commConsts.Run, wsMsg)
			}

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

func ExeScript(scriptFile, workspacePath string, conf commDomain.WorkspaceConf, report *commDomain.ZtfReport, scriptIdx,
	total, pathMaxWidth, numbMaxWidth int,
	ch chan int, wsMsg *websocket.Message) {

	startTime := time.Now()

	startMsg := i118Utils.Sprintf("start_execution", scriptFile, dateUtils.DateTimeStr(startTime))

	if commConsts.ExecFrom != commConsts.FromCmd {
		websocketHelper.SendExecMsg(startMsg, "", commConsts.Run, wsMsg)
		logUtils.ExecConsolef(-1, startMsg)
	}

	logUtils.ExecFilef(startMsg)

	logs := ""
	stdOutput, errOutput := RunScript(scriptFile, workspacePath, conf, ch, wsMsg)
	stdOutput = strings.Trim(stdOutput, "\n")

	if stdOutput != "" {
		logs = stdOutput
	}
	if errOutput != "" {
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendOutputMsg(errOutput, "", wsMsg)
		}
		logUtils.ExecConsolef(-1, errOutput)
		logUtils.ExecFilef(errOutput)
	}

	entTime := time.Now()
	secs := fmt.Sprintf("%.2f", float32(entTime.Sub(startTime)/time.Second))

	endMsg := i118Utils.Sprintf("end_execution", scriptFile, dateUtils.DateTimeStr(entTime))
	if commConsts.ExecFrom != commConsts.FromCmd {
		websocketHelper.SendExecMsg(endMsg, "", commConsts.Run, wsMsg)
		logUtils.ExecConsolef(-1, endMsg)
	}

	logUtils.ExecFilef(endMsg)

	CheckCaseResult(scriptFile, logs, report, scriptIdx, total, secs, pathMaxWidth, numbMaxWidth, wsMsg)
}

func RunScript(filePath, workspacePath string, conf commDomain.WorkspaceConf,
	ch chan int, wsMsg *websocket.Message) (
	stdOutput string, errOutput string) {

	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		lang := langHelper.GetLangByFile(filePath)

		scriptInterpreter := ""
		if strings.ToLower(lang) != "bat" {
			scriptInterpreter = configHelper.GetFieldVal(conf, stringUtils.UcFirst(lang))
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
			if commConsts.ExecFrom != commConsts.FromCmd {
				websocketHelper.SendOutputMsg(msg, "", wsMsg)
			}
			logUtils.ExecConsolef(-1, msg)
			logUtils.ExecFilef(msg)
		}
	} else {
		err := os.Chmod(filePath, 0777)
		if err != nil {
			msg := i118Utils.I118Prt.Sprintf("exec_cmd_fail", filePath, err.Error())
			if commConsts.ExecFrom != commConsts.FromCmd {
				websocketHelper.SendOutputMsg(msg, "", wsMsg)
			}
			logUtils.ExecConsolef(-1, msg)
			logUtils.ExecFilef(msg)
		}

		filePath = "\"" + filePath + "\""
		cmd = exec.Command("/bin/bash", "-c", filePath)
	}

	cmd.Dir = workspacePath

	if cmd == nil {
		msgStr := i118Utils.Sprintf("cmd_empty")
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendOutputMsg(msgStr, "", wsMsg)
			logUtils.ExecConsolef(color.FgRed, msgStr)
		}

		logUtils.ExecFilef(msgStr)

		return "", msgStr
	}

	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendOutputMsg(err1.Error(), "", wsMsg)
		}
		logUtils.ExecConsolef(color.FgRed, err1.Error())
		logUtils.ExecFilef(err1.Error())

		return "", err1.Error()
	} else if err2 != nil {
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendOutputMsg(err2.Error(), "", wsMsg)
		}
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

			if commConsts.ExecFrom != commConsts.FromCmd {
				websocketHelper.SendOutputMsg(line, "", wsMsg)
				logUtils.ExecConsole(1, line)
			}

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

			if commConsts.ExecFrom != commConsts.FromCmd {
				websocketHelper.SendExecMsg(msg, "", commConsts.Run, wsMsg)
			}

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

func FilterCases(cases []string, conf commDomain.WorkspaceConf) (casesToRun, casesToIgnore []string) {
	for _, cs := range cases {
		if commonUtils.IsWin() {
			if path.Ext(cs) == ".sh" { // filter by os
				continue
			}

			ext := path.Ext(cs)
			if ext != "" {
				ext = ext[1:]
			}
			lang := commConsts.ScriptExtToNameMap[ext]
			interpreter := configHelper.GetFieldVal(conf, stringUtils.UcFirst(lang))
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

func genReport(productId, id int, by commConsts.ExecBy) (report commDomain.ZtfReport) {
	report = commDomain.ZtfReport{
		TestEnv: commonUtils.GetOs(), ExecBy: by, ExecById: id, ProductId: productId,
		Pass: 0, Fail: 0, Total: 0, FuncResult: make([]commDomain.FuncResult, 0)}

	report.TestType = commConsts.TestFunc
	report.TestTool = commConsts.AppServer

	return
}

func getNumbMaxWidth(casesToRun []string) (numbMaxWidth, pathMaxWidth int) {
	for _, cs := range casesToRun {
		lent := runewidth.StringWidth(cs)
		if lent > pathMaxWidth {
			pathMaxWidth = lent
		}

		content := fileUtils.ReadFile(cs)
		caseId := scriptHelper.ReadCaseId(content)
		if len(caseId) > numbMaxWidth {
			numbMaxWidth = len(caseId)
		}
	}

	return
}
