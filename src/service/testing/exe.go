package testingService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/shell"
	"path"
	"strings"
	"time"
)

func ExeScripts(files []string, report *model.TestReport, pathMaxWidth int) {
	casesToRun := make([]string, 0)

	for _, file := range files {
		if (commonUtils.IsWin() && path.Ext(file) != ".sh") || (!commonUtils.IsWin() && path.Ext(file) != ".bat") {
			casesToRun = append(casesToRun, file)
		}
	}

	now := time.Now()
	startTime := now.Unix()
	report.StartTime = startTime

	logUtils.ScreenAndResult(now.Format("2006-01-02 15:04:05") + " " +
		i118Utils.I118Prt.Sprintf("found_scripts", len(casesToRun)) + "\n")

	for idx, file := range casesToRun {
		ExeScript(file, report, idx, len(casesToRun), pathMaxWidth)
	}

	endTime := time.Now().Unix()
	secs := endTime - startTime

	report.EndTime = startTime
	report.Duration = secs
}

func ExeScript(file string, report *model.TestReport, idx int, total int, pathMaxWidth int) {
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

	CheckCaseResult(file, logs, report, idx, total, secs, pathMaxWidth)

	if idx < total-1 {
		logUtils.Log("")
	}
}
