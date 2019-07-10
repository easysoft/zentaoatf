package biz

import (
	"os"
	"runtime"
	"strings"
	"utils"
)

func RunScripts(files []string, dir string, langType string) {
	logDir := dir + string(os.PathSeparator) + "logs"
	utils.MkDir(logDir)

	for _, file := range files {
		RunScript(file, langType)
	}
}

func RunScript(file string, langType string) {
	osName := runtime.GOOS

	var command string
	var logFile string
	if osName == "darwin" {
		logFile = utils.ScriptToLog(file)
		command = file //  + " > " + logFile

		if langType == "php" {
			command = langType + " " + command
		}
	}

	output := utils.ExecCommand(command)
	utils.WriteFile(logFile, strings.Join(output, ""))
}
