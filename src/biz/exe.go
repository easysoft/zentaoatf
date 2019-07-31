package biz

import (
	"github.com/easysoft/zentaoatf/src/misc"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/fatih/color"
	"strings"
	"time"
)

func ExeScripts(files []string, dir string, langType string, report *model.TestReport) {
	PrintWholeLine(utils.I118Prt.Sprintf("start_execution", ""), "=", color.FgCyan)

	startTime := time.Now().Unix()
	report.StartTime = startTime

	for _, file := range files {
		ExeScript(file, langType, dir)
	}

	PrintWholeLine(utils.I118Prt.Sprintf("end_execution", ""), "=", color.FgCyan)

	endTime := time.Now().Unix()
	secs := endTime - startTime

	report.EndTime = startTime
	report.Duration = secs
}

func ExeScript(file string, langType string, dir string) {
	var command string
	var logFile string

	if utils.IsMac() {
		logFile = utils.ScriptToLogName(dir, file)
		command = file //  + " > " + logFile

		if langType == misc.PHP.String() {
			command = langType + " " + command
		}
	}

	startTime := time.Now()

	msg := utils.I118Prt.Sprintf("start_case", file, startTime.Format("2006-01-02 15:04:05"))
	PrintWholeLine(msg, "-", color.FgCyan)

	output := utils.ExecCommand(command)
	utils.WriteFile(logFile, strings.Join(output, ""))
	if utils.RunFromCui {
		Printt(strings.Join(output, ""))
	}

	entTime := time.Now()
	secs := int64(entTime.Sub(startTime) / time.Second)

	msg = utils.I118Prt.Sprintf("end_case", file, entTime.Format("2006-01-02 15:04:05"), secs)
	PrintWholeLine(msg, "-", color.FgCyan)
}
