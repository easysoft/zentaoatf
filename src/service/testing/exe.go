package testingService

import (
	"github.com/easysoft/zentaoatf/src/model"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/shell"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/fatih/color"
	"time"
)

func ExeScripts(files []string, report *model.TestReport) {
	logUtils.PrintWholeLine(i118Utils.I118Prt.Sprintf("start_execution", ""), "=", color.FgCyan)

	startTime := time.Now().Unix()
	report.StartTime = startTime

	for _, file := range files {
		ExeScript(file, report)
	}

	logUtils.PrintWholeLine(i118Utils.I118Prt.Sprintf("end_execution", ""), "=", color.FgCyan)

	endTime := time.Now().Unix()
	secs := endTime - startTime

	report.EndTime = startTime
	report.Duration = secs
}

func ExeScript(file string, report *model.TestReport) {
	var logFile string

	logFile = zentaoUtils.ScriptToLogName(file)
	logUtils.InitLog(zentaoUtils.ScriptToLogDir(file))

	startTime := time.Now()

	msg := i118Utils.I118Prt.Sprintf("start_case", file, startTime.Format("2006-01-02 15:04:05"))
	logUtils.PrintWholeLine(msg, "-", color.FgCyan)

	logUtils.Screen(msg)
	logUtils.Trace(msg)

	output := shellUtils.ExecFile(file)
	fileUtils.WriteFile(logFile, output)

	CheckResult(file, report)

	entTime := time.Now()
	secs := int64(entTime.Sub(startTime) / time.Second)

	msg = i118Utils.I118Prt.Sprintf("end_case", file, entTime.Format("2006-01-02 15:04:05"), secs)
	logUtils.PrintWholeLine(msg, "-", color.FgCyan)
}
