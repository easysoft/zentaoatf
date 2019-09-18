package testingService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/shell"
	"strings"
	"time"
)

func ExeScripts(files []string, report *model.TestReport, pathMaxWidth int) {
	//msg := i118Utils.I118Prt.Sprint("start_execution")
	//msg = logUtils.GetWholeLine(msg, "=")
	//logUtils.Trace(msg)

	now := time.Now()
	startTime := now.Unix()
	report.StartTime = startTime

	logUtils.ScreenAndResult(now.Format("2006-01-02 15:04:05") + " " +
		i118Utils.I118Prt.Sprintf("found_scripts", len(files)) + "\n")

	for idx, file := range files {
		ExeScript(file, report, idx, len(files), pathMaxWidth)
	}

	//msg = i118Utils.I118Prt.Sprint("end_execution")
	//logUtils.Trace(logUtils.GetWholeLine(msg, "=") + "\n")

	endTime := time.Now().Unix()
	secs := endTime - startTime

	report.EndTime = startTime
	report.Duration = secs
}

func ExeScript(file string, report *model.TestReport, idx int, total int, pathMaxWidth int) {
	startTime := time.Now()

	logUtils.Log("===start " + file + " at " + startTime.Format("2006-01-02 15:04:05"))

	output := shellUtils.ExecFile(file)
	output = strings.Trim(output, "\n")

	if output != "" {
		logUtils.Log(output)
	}

	entTime := time.Now()
	secs := fmt.Sprintf("%.2f", float32(entTime.Sub(startTime)/time.Second))

	CheckCaseResult(file, report, idx, total, secs, pathMaxWidth)

	logUtils.Log("===end " + file + " at " + entTime.Format("2006-01-02 15:04:05"))

	if idx < total-1 {
		logUtils.Log("")
	}
}
