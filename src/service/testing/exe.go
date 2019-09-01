package testingService

import (
	"github.com/easysoft/zentaoatf/src/model"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/shell"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"time"
)

func ExeScripts(files []string, report *model.TestReport) {
	logUtils.InitLog(zentaoUtils.ScriptToLogDir())

	msg := i118Utils.I118Prt.Sprintf("start_execution", "")
	msg = logUtils.GetWholeLine(msg, "=")
	logUtils.Trace(msg)

	startTime := time.Now().Unix()
	report.StartTime = startTime

	for idx, file := range files {
		ExeScript(file, report, idx, len(files))
	}

	msg = i118Utils.I118Prt.Sprintf("end_execution", "")
	logUtils.Trace(logUtils.GetWholeLine(msg, "=") + "\n")

	endTime := time.Now().Unix()
	secs := endTime - startTime

	report.EndTime = startTime
	report.Duration = secs
}

func ExeScript(file string, report *model.TestReport, idx int, total int) {
	logFile := zentaoUtils.ScriptToLogName(file)

	startTime := time.Now()

	msg := i118Utils.I118Prt.Sprintf("start_case", file, startTime.Format("2006-01-02 15:04:05"))
	logUtils.Trace(logUtils.GetWholeLine(msg, "-"))

	output := shellUtils.ExecFile(file)
	fileUtils.WriteFile(logFile, output)

	CheckCaseResult(file, report, idx, total)

	entTime := time.Now()
	secs := int64(entTime.Sub(startTime) / time.Second)

	msg = i118Utils.I118Prt.Sprintf("end_case", file, entTime.Format("2006-01-02 15:04:05"), secs)
	logUtils.Trace(logUtils.GetWholeLine(msg, "-"))

	if idx < total-1 {
		logUtils.Trace("")
	}
}
