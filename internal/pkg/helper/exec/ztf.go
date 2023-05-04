package execHelper

import (
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	analysisHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/analysis"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	channelUtils "github.com/easysoft/zentaoatf/pkg/lib/channel"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"github.com/mattn/go-runewidth"
)

func ExecCases(ch chan int, testSet serverDomain.TestSet, msg *websocket.Message) (report commDomain.ZtfReport, pathMaxWidth int, err error) {
	// if testSet.Cmd != "" {
	// 	err = execCmd(testSet.Cmd, testSet.WorkspacePath)
	// 	if err != nil {
	// 		return
	// 	}
	// }

	return RunZtf(ch, testSet.WorkspacePath, 0, 0, commConsts.Case, testSet.Cases, msg)
}

func execCmd(cmd string, workspacePath string) (err error) {
	pth := filepath.Join(workspacePath, ".cmd.tmp")
	fileUtils.WriteFile(pth, cmd)

	conf := configHelper.LoadByWorkspacePath(workspacePath)

	stdOutput, errOutput := RunFile(pth, workspacePath, conf, nil, nil, 0)
	if errOutput != "" {
		logUtils.Infof("failed to exec command '%s' without output %s, err %v.", pth, stdOutput, errOutput)
	} else {
		logUtils.Infof("exec command '%s' with output %v.", pth, stdOutput)
	}

	return
}

func ExecModule(ch chan int, testSet serverDomain.TestSet, msg *websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {

	cases, err := zentaoHelper.GetCasesByModuleInDir(testSet.ProductId, testSet.ModuleId,
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

	casesToRun, casesToIgnore := FilterCases(cases, &conf)

	numbMaxWidth := 0
	numbMaxWidth, pathMaxWidth = getNumbMaxWidth(casesToRun)
	report = genReport(productId, id, by)

	// exec scripts
	ExeScripts(casesToRun, casesToIgnore, workspacePath, conf, &report, pathMaxWidth, numbMaxWidth, ch, wsMsg)

	// gen report
	if len(casesToRun) > 0 {
		GenZTFTestReport(report, pathMaxWidth, workspacePath, ch, wsMsg)
	}

	if commConsts.ExecFrom == commConsts.FromClient {
		websocketHelper.SendExecMsg("", "false", commConsts.Run, nil, wsMsg)
	}
	if ch != nil {
		if !channelUtils.IsChanClose(ch) {
			close(ch)
		}
	}

	return
}

func ExeScripts(casesToRun []string, casesToIgnore []string, workspacePath string, conf commDomain.WorkspaceConf,
	report *commDomain.ZtfReport, pathMaxWidth int, numbMaxWidth int,
	ch chan int, wsMsg *websocket.Message) {

	now := time.Now()
	startTime := now.Unix()
	report.StartTime = startTime

	workDir := commConsts.WorkDir
	if commConsts.ExecFrom == commConsts.FromClient {
		workDir = workspacePath
	}

	msg := i118Utils.Sprintf("found_scripts", len(casesToRun), workDir, commConsts.ZtfDir)

	if commConsts.ExecFrom == commConsts.FromClient {
		msg = i118Utils.Sprintf("found_scripts_no_ztf_dir", len(casesToRun), workDir)
		websocketHelper.SendExecMsg(msg, "", commConsts.Run, nil, wsMsg)
	}
	logUtils.ExecConsolef(color.FgCyan, msg)
	logUtils.ExecResult(msg)

	if len(casesToIgnore) > 0 {
		temp := i118Utils.Sprintf("ignore_scripts", strconv.Itoa(len(casesToIgnore)))
		if commConsts.ExecFrom == commConsts.FromClient {
			websocketHelper.SendExecMsg(temp, "", commConsts.Run, nil, wsMsg)
		}

		logUtils.ExecConsolef(color.FgCyan, temp)
		logUtils.ExecResult(temp)
	}

	if commConsts.BatchCount <= 1 {
		execScripts(casesToRun, workspacePath, conf, report, pathMaxWidth, numbMaxWidth, ch, wsMsg)
	} else {
		batchExecScripts(casesToRun, workspacePath, conf, report, pathMaxWidth, numbMaxWidth, ch, wsMsg)
	}
	endTime := time.Now().Unix()
	report.EndTime = endTime
	report.Duration = endTime - startTime
}

func execScripts(casesToRun []string, workspacePath string, conf commDomain.WorkspaceConf,
	report *commDomain.ZtfReport, pathMaxWidth int, numbMaxWidth int,
	ch chan int, wsMsg *websocket.Message) {
	for idx, file := range casesToRun {
		if fileUtils.IsDir(file) {
			continue
		}

		ExecScript(file, workspacePath, conf, report, idx, len(casesToRun), pathMaxWidth, numbMaxWidth, ch, wsMsg, nil)

		select {
		case <-ch:
			msg := i118Utils.Sprintf("exit_exec_all")
			if commConsts.ExecFrom == commConsts.FromClient {
				websocketHelper.SendExecMsg(msg, "", commConsts.Run, nil, wsMsg)
			}

			logUtils.ExecConsolef(color.FgCyan, msg)
			logUtils.ExecFilef(msg)

			return
		default:
		}
	}
}

func batchExecScripts(casesToRun []string, workspacePath string, conf commDomain.WorkspaceConf,
	report *commDomain.ZtfReport, pathMaxWidth int, numbMaxWidth int,
	ch chan int, wsMsg *websocket.Message) {
	realCasesToRun := []string{}
	for _, file := range casesToRun {
		if fileUtils.IsDir(file) {
			continue
		}
		realCasesToRun = append(realCasesToRun, file)
	}

	var wg sync.WaitGroup
	var lock sync.Mutex
	poolId := 0
	caseCount := len(realCasesToRun)

	wgCount := commConsts.BatchCount
	if caseCount-poolId < commConsts.BatchCount {
		wgCount = caseCount - poolId
	}
	wg.Add(wgCount)

	for _, file := range realCasesToRun {
		poolId++
		go func(poolId int, file string) {
			ExecScript(file, workspacePath, conf, report, poolId-1, caseCount, pathMaxWidth, numbMaxWidth, ch, wsMsg, &lock)
			wg.Done()
		}(poolId, file)

		if caseCount > commConsts.BatchCount && poolId%commConsts.BatchCount == 0 {
			wg.Wait()
			if caseCount-poolId < commConsts.BatchCount {
				wgCount = caseCount - poolId
			}
			wg.Add(wgCount)
		}

		select {
		case <-ch:
			msg := i118Utils.Sprintf("exit_exec_all")
			if commConsts.ExecFrom == commConsts.FromClient {
				websocketHelper.SendExecMsg(msg, "", commConsts.Run, nil, wsMsg)
			}

			logUtils.ExecConsolef(color.FgCyan, msg)
			logUtils.ExecFilef(msg)

			return
		default:
		}
	}

	if caseCount <= commConsts.BatchCount {
		wg.Wait()
	}
}

func FilterCases(cases []string, conf *commDomain.WorkspaceConf) (casesToRun, casesToIgnore []string) {
	for _, cs := range cases {
		ext := path.Ext(cs)
		if ext != "" {
			ext = ext[1:]
		}
		lang := commConsts.ScriptExtToNameMap[ext]
		if lang == "" {
			continue
		}

		if commonUtils.IsWin() {
			if path.Ext(cs) == ".sh" { // filter by os
				continue
			}

			interpreter := configHelper.GetFieldVal(*conf, stringUtils.UcFirst(lang))
			if interpreter == "-" || interpreter == "" {
				interpreter = ""
				if lang != "bat" {
					ok := AddInterpreterIfExist(conf, lang)
					if !ok {
						casesToIgnore = append(casesToIgnore, cs)
					} else {
						interpreter = configHelper.GetFieldVal(*conf, stringUtils.UcFirst(lang))
					}
				}
			}
			if lang != "bat" && interpreter == "" { // ignore the ones with no interpreter set
				continue
			}

		} else { // filter by os
			if path.Ext(cs) == ".bat" {
				continue
			}
		}

		//pass := scriptHelper.CheckFileIsScript(cs)
		//if pass {
		casesToRun = append(casesToRun, cs)
		//}
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
