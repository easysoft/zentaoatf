package testingService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/shell"
	"github.com/fatih/color"
	"strconv"
	"strings"
	"time"
)

func ExeScripts(casesToRun []string, casesToIgnore []string, report *model.TestReport, pathMaxWidth int, numbMaxWidth int) {
	now := time.Now()
	startTime := now.Unix()
	report.StartTime = startTime

	postFix := ":"
	if len(casesToRun) == 0 {
		postFix = "."
	}

	logUtils.Result("\n" + logUtils.GetWholeLine(now.Format("2006-01-02 15:04:05")+" "+
		i118Utils.Sprintf("found_scripts", strconv.Itoa(len(casesToRun)))+postFix, "="))
	logUtils.Screen("\n" + logUtils.GetWholeLine(now.Format("2006-01-02 15:04:05")+" "+
		i118Utils.Sprintf("found_scripts", color.CyanString(strconv.Itoa(len(casesToRun))))+postFix, "="))

	if len(casesToIgnore) > 0 {
		logUtils.Result("                    " +
			i118Utils.Sprintf("ignore_scripts", strconv.Itoa(len(casesToIgnore))) + postFix)
		logUtils.Screen("                    " +
			i118Utils.Sprintf("ignore_scripts", color.CyanString(strconv.Itoa(len(casesToIgnore)))) + postFix)
	}

	for idx, file := range casesToRun {
		ExeScript(file, report, idx, len(casesToRun), pathMaxWidth, numbMaxWidth)
	}

	endTime := time.Now().Unix()
	report.EndTime = endTime
	report.Duration = endTime - startTime
}

func ExeScript(file string, report *model.TestReport, idx int, total int, pathMaxWidth int, numbMaxWidth int) {
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
