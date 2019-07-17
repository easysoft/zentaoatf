package biz

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/misc"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/fatih/color"
	"strings"
	"time"
)

func ExeScripts(files []string, dir string, langType string, report *model.TestReport) {
	fmt.Println(color.BlueString("=== Begin to run test scripts ==="))

	startTime := time.Now().Unix()
	report.StartTime = startTime

	for _, file := range files {
		ExeScript(file, langType, dir)
	}

	fmt.Println(color.BlueString("=== End to run test scripts ==="))

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
	fmt.Printf(color.CyanString("\n--- Start %s %s"), file, startTime.Format("2006-01-02 15:04:05")+"\n")

	output := utils.ExecCommand(command)
	utils.WriteFile(logFile, strings.Join(output, ""))

	entTime := time.Now()
	secs := int64(entTime.Sub(startTime) / time.Second)
	fmt.Printf(color.CyanString("--- End %s %dsec %s"), file, secs, "\n")
}
