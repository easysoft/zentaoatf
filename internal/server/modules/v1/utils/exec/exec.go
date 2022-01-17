package scriptUtils

import (
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
	"github.com/mattn/go-runewidth"
	"path"
	"strconv"
	"strings"
	"time"
)

func ExecCase(req serverDomain.TestExec, projectPath string) (report commDomain.ZtfReport, pathMaxWidth int, err error) {
	casesToRun, casesToIgnore := filterCases(req.Cases, projectPath)

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

	ExeScripts(casesToRun, casesToIgnore, &report, pathMaxWidth, numbMaxWidth)

	return
}

func filterCases(cases []string, projectPath string) (casesToRun, casesToIgnore []string) {
	conf := configUtils.LoadByProjectPath(projectPath)

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

func ExeScripts(casesToRun []string, casesToIgnore []string, report *commDomain.ZtfReport, pathMaxWidth int, numbMaxWidth int) {
	now := time.Now()
	startTime := now.Unix()
	report.StartTime = startTime

	postFix := ":"
	if len(casesToRun) == 0 {
		postFix = "."
	}

	logUtils.Result(i118Utils.Sprintf("found_scripts", strconv.Itoa(len(casesToRun)))+postFix, "="))

	if len(casesToIgnore) > 0 {
		logUtils.Result(i118Utils.Sprintf("ignore_scripts", strconv.Itoa(len(casesToIgnore))) + postFix)
	}

	for idx, file := range casesToRun {
		ExeScript(file, report, idx, len(casesToRun), pathMaxWidth, numbMaxWidth)
	}

	endTime := time.Now().Unix()
	report.EndTime = endTime
	report.Duration = endTime - startTime
}

func ExeScript(file string, report *commDomain.ZtfReport, idx int, total int, pathMaxWidth int, numbMaxWidth int) {
	startTime := time.Now()

	logUtils.Log("===start " + file + " at " + startTime.Format("2006-01-02 15:04:05"))
	logs := ""

	out, err := shellUtils.ExecScriptFile(file)
	out = strings.Trim(out, "\n")

	if out != "" {
		logUtils.Log(out)
		logs = out
	}
	if err != "" {
		logUtils.Error(err)
	}

	entTime := time.Now()
	secs := fmt.Sprintf("%.2f", float32(entTime.Sub(startTime)/time.Second))

	logUtils.Log("===end " + file + " at " + entTime.Format("2006-01-02 15:04:05"))
	CheckCaseResult(file, logs, report, idx, total, secs, pathMaxWidth, numbMaxWidth)

	if idx < total-1 {
		logUtils.Log("")
	}
}