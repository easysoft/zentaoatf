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

	postFix := ":\n"
	if len(casesToRun) == 0 {
		postFix = "."
	}

	logUtils.ScreenAndResult(now.Format("2006-01-02 15:04:05") + " " +
		i118Utils.I118Prt.Sprintf("found_scripts", color.CyanString(strconv.Itoa(len(casesToRun)))) + postFix)

	if len(casesToIgnore) > 0 {
		logUtils.ScreenAndResult("                    " +
			i118Utils.I118Prt.Sprintf("ignore_scripts", color.CyanString(strconv.Itoa(len(casesToIgnore)))) + postFix)
	}

	for idx, file := range casesToRun {
		ExeScript(file, report, idx, len(casesToRun), pathMaxWidth, numbMaxWidth)
	}

	endTime := time.Now().Unix()
	secs := endTime - startTime

	report.EndTime = startTime
	report.Duration = secs
}

func ExeScript(file string, report *model.TestReport, idx int, total int, pathMaxWidth int, numbMaxWidth int) {
	startTime := time.Now()

	logUtils.Log("===start " + file + " at " + startTime.Format("2006-01-02 15:04:05"))
	logs := ""

	output := shellUtils.ExecFile(file)
	output = strings.Trim(output, "\n")

	if output != "" {
		logUtils.Log(output)
		logs = output
	}

	entTime := time.Now()
	secs := fmt.Sprintf("%.2f", float32(entTime.Sub(startTime)/time.Second))

	logUtils.Log("===end " + file + " at " + entTime.Format("2006-01-02 15:04:05"))

	CheckCaseResult(file, logs, report, idx, total, secs, pathMaxWidth, numbMaxWidth)

	if idx < total-1 {
		logUtils.Log("")
	}
}
