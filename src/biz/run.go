package biz

import (
	"fmt"
	"model"
	"strings"
	"time"
	"utils"
)

func RunScripts(files []string, dir string, langType string, report *model.TestReport) {
	fmt.Println("=== Begin to run test scripts ===")

	startTime := time.Now().Unix()
	report.StartTime = startTime

	for _, file := range files {
		RunScript(file, langType, dir)
	}

	endTime := time.Now().Unix()
	secs := endTime - startTime

	report.EndTime = startTime
	report.Duration = secs
}

func RunScript(file string, langType string, dir string) {
	var command string
	var logFile string

	if utils.IsMac() {
		logFile = utils.ScriptToLogName(dir, file)
		command = file //  + " > " + logFile

		if langType == "php" {
			command = langType + " " + command
		}
	}

	startTime := time.Now()
	fmt.Printf("\n--- Start %s %s", file, startTime.Format("2006-01-02 15:04:05")+"\n")

	output := utils.ExecCommand(command)
	utils.WriteFile(logFile, strings.Join(output, ""))

	entTime := time.Now()
	secs := int64(entTime.Sub(startTime) / time.Second)
	fmt.Printf("--- End %s %dsec %s", file, secs, "\n")
}
