package biz

import (
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/fatih/color"
	"log"
	"time"
)

func ExeScripts(files []string, scriptDir string, langType string, report *model.TestReport) {
	utils.PrintWholeLine(utils.I118Prt.Sprintf("start_execution", ""), "=", color.FgCyan)

	startTime := time.Now().Unix()
	report.StartTime = startTime

	for _, file := range files {
		ExeScript(file, langType, scriptDir)
	}

	utils.PrintWholeLine(utils.I118Prt.Sprintf("end_execution", ""), "=", color.FgCyan)

	endTime := time.Now().Unix()
	secs := endTime - startTime

	report.EndTime = startTime
	report.Duration = secs
}

func ExeScript(file string, langType string, scriptDir string) {
	var command string
	var logFile string

	logFile = utils.ScriptToLogName(scriptDir, file)
	log.Panic(logFile)
	command = file

	startTime := time.Now()

	msg := utils.I118Prt.Sprintf("start_case", file, startTime.Format("2006-01-02 15:04:05"))
	utils.PrintWholeLine(msg, "-", color.FgCyan)

	output := utils.ExecFile(command)
	utils.WriteFile(logFile, output)

	entTime := time.Now()
	secs := int64(entTime.Sub(startTime) / time.Second)

	msg = utils.I118Prt.Sprintf("end_case", file, entTime.Format("2006-01-02 15:04:05"), secs)
	utils.PrintWholeLine(msg, "-", color.FgCyan)
}
