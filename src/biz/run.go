package biz

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
	"utils"
)

func RunScripts(files []string, dir string, langType string, summaryMap *map[string]interface{}) {
	logDir := dir + string(os.PathSeparator) + "logs"
	utils.MkDir(logDir)

	fmt.Println("=== Begin to run test scripts ===")

	startTime := time.Now().Unix()
	(*summaryMap)["startTime"] = startTime

	for _, file := range files {
		RunScript(file, langType)
	}

	endTime := time.Now().Unix()
	secs := endTime - startTime
	(*summaryMap)["endTime"] = startTime
	(*summaryMap)["duration"] = secs
}

func RunScript(file string, langType string) {
	osName := runtime.GOOS

	var command string
	var logFile string
	if osName == "darwin" {
		logFile = utils.ScriptToLogName(file)
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
