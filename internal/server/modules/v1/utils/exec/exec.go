package scriptUtils

import (
	"bufio"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/zentao"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/config"
	resultUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/result"
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

func Exec(ch chan int, fun func(info string, msg websocket.Message), req serverDomain.WsMsg, msg websocket.Message) (
	err error) {
	if req.Act == commConsts.ExecCase {
		ExecCase(ch, fun, req, msg)
	}

	return
}

func ExecCase(ch chan int, fun func(info string, msg websocket.Message), req serverDomain.WsMsg, msg websocket.Message) (report commDomain.ZtfReport, pathMaxWidth int, err error) {

	projectPath := req.ProjectPath
	conf := configUtils.LoadByProjectPath(projectPath)

	casesToRun, casesToIgnore := filterCases(req.Cases, conf)

	report = commDomain.ZtfReport{Env: commonUtils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, FuncResult: make([]commDomain.FuncResult, 0)}
	report.TestType = "func"
	report.TestFrame = commConsts.AppServer

	numbMaxWidth := 0
	for _, cs := range casesToRun {
		lent := runewidth.StringWidth(cs)
		if lent > pathMaxWidth {
			pathMaxWidth = lent
		}

		content := fileUtils.ReadFile(cs)
		caseId := zentaoUtils.ReadCaseId(content)
		if len(caseId) > numbMaxWidth {
			numbMaxWidth = len(caseId)
		}
	}

	ExeScripts(casesToRun, casesToIgnore, projectPath, conf, &report, pathMaxWidth, numbMaxWidth, ch, fun, msg)

	return
}

func ExeScripts(casesToRun []string, casesToIgnore []string, projectPath string, conf commDomain.ProjectConf,
	report *commDomain.ZtfReport, pathMaxWidth int, numbMaxWidth int, ch chan int, fun func(info string, msg websocket.Message), msg websocket.Message) {

	now := time.Now()
	startTime := now.Unix()
	report.StartTime = startTime

	postFix := ":"
	if len(casesToRun) == 0 {
		postFix = "."
	}

	temp := i118Utils.Sprintf("found_scripts", strconv.Itoa(len(casesToRun))) + postFix
	fun(temp, msg)

	if len(casesToIgnore) > 0 {
		temp := i118Utils.Sprintf("ignore_scripts", strconv.Itoa(len(casesToIgnore))) + postFix
		fun(temp, msg)
	}

	for idx, file := range casesToRun {
		ExeScript(file, projectPath, conf, report, idx, len(casesToRun), pathMaxWidth, numbMaxWidth, ch, fun, msg)
	}

	endTime := time.Now().Unix()
	report.EndTime = endTime
	report.Duration = endTime - startTime
}

func ExeScript(file, projectPath string, conf commDomain.ProjectConf, report *commDomain.ZtfReport, idx,
	total, pathMaxWidth, numbMaxWidth int, ch chan int, fun func(s string, msg websocket.Message), msg websocket.Message) {

	startTime := time.Now()

	fun("===start "+file+" at "+startTime.Format("2006-01-02 15:04:05"), msg)
	logs := ""

	out, err := ExecScriptFile(file, projectPath, conf, ch, fun, msg)
	out = strings.Trim(out, "\n")

	if out != "" {
		logUtils.Infof(out)
		logs = out
	}
	if err != "" {
		logUtils.Error(err)
	}

	entTime := time.Now()
	secs := fmt.Sprintf("%.2f", float32(entTime.Sub(startTime)/time.Second))

	logUtils.Infof("===end " + file + " at " + entTime.Format("2006-01-02 15:04:05"))
	resultUtils.CheckCaseResult(file, logs, report, idx, total, secs, pathMaxWidth, numbMaxWidth)

	if idx < total-1 {
		logUtils.Infof("")
	}
}

func ExecScriptFile(filePath, projectPath string, conf commDomain.ProjectConf,
	ch chan int, fun func(s string, msg websocket.Message), msg websocket.Message) (ret1 string, ret2 string) {
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
			fmt.Printf("use interpreter %s for script %s\n", scriptInterpreter, filePath)
			i118Utils.I118Prt.Printf("no_interpreter_for_run", filePath, lang)
		}
	} else {
		err := os.Chmod(filePath, 0777)
		if err != nil {
			logUtils.Infof("chmod error" + err.Error())
		}

		filePath = "\"" + filePath + "\""
		cmd = exec.Command("/bin/bash", "-c", filePath)
	}

	cmd.Dir = projectPath

	if cmd == nil {
		msg := "error cmd is nil"
		logUtils.Infof(msg)
		return "", fmt.Sprint(msg)
	}

	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		fmt.Println(err1)
		return "", fmt.Sprint(err1)
	} else if err2 != nil {
		fmt.Println(err2)
		return "", fmt.Sprint(err2)
	}

	cmd.Start()

	reader1 := bufio.NewReader(stdout)
	output1 := make([]string, 0)
	for {
		line, err2 := reader1.ReadString('\n')
		if line != "" {
			fun(line, msg)
		}

		if err2 != nil || io.EOF == err2 {
			break
		}
		output1 = append(output1, line)

		select {
		case <-ch:
			fmt.Println("exiting...")
			ch <- 1
			break
		default:
			fmt.Println("continue...")
		}
	}

	reader2 := bufio.NewReader(stderr)
	output2 := make([]string, 0)
	for {
		line, err2 := reader2.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		output2 = append(output2, line)
	}

	cmd.Wait()

	ret1 = strings.Join(output1, "")
	ret2 = strings.Join(output2, "")
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
