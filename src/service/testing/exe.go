package testingService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/shell"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"time"
	"unicode/utf8"
)

func ExeScripts(files []string, report *model.TestReport) {
	logUtils.InitLog(zentaoUtils.ScriptToLogDir())

	//msg := i118Utils.I118Prt.Sprint("start_execution")
	//msg = logUtils.GetWholeLine(msg, "=")
	//logUtils.Trace(msg)

	startTime := time.Now().Unix()
	report.StartTime = startTime

	pathMaxWidth := 0
	for _, file := range files {
		lent := utf8.RuneCountInString(file)
		if lent > pathMaxWidth {
			pathMaxWidth = lent
		}
	}

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
	logFile := zentaoUtils.ScriptToLogName(file)

	startTime := time.Now()

	//msg := i118Utils.I118Prt.Sprintf("start_case", file, startTime.Format("2006-01-02 15:04:05"))
	//logUtils.Trace(logUtils.GetWholeLine(msg, "-"))

	output := shellUtils.ExecFile(file)
	fileUtils.WriteFile(logFile, output)

	entTime := time.Now()
	secs := fmt.Sprintf("%.2f", float32(entTime.Sub(startTime)/time.Second))

	CheckCaseResult(file, report, idx, total, secs, pathMaxWidth)

	//msg = i118Utils.I118Prt.Sprintf("end_case", file, entTime.Format("2006-01-02 15:04:05"), secs)
	//logUtils.Trace(logUtils.GetWholeLine(msg, "-"))
	//
	//if idx < total-1 {
	//	logUtils.Trace("")
	//}
}
