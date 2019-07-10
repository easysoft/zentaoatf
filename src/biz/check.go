package biz

import (
	"fmt"
	"regexp"
	"strings"
	"utils"
)

func CheckResults(dir string, langType string) {
	scriptFiles, _ := utils.GetAllFiles(dir, langType)

	for _, scriptFile := range scriptFiles {
		logFile := utils.ScriptToLog(scriptFile)

		expectContent := utils.ReadExpect(scriptFile)
		logContent := utils.ReadFile(logFile)

		expectContent = strings.Trim(expectContent, "\n")
		logContent = strings.Trim(logContent, "\n")

		fmt.Printf(scriptFile + ": " + logFile + "\n")
		Compare(expectContent, logContent)
	}
}

func Compare(expectContent string, logContent string) {
	expectArr := strings.Split(expectContent, "\n")
	logArr := strings.Split(logContent, "\n")

	fmt.Printf("%d %d \n", len(expectArr), len(logArr))

	for numb, line := range expectArr {
		log := "N/A"
		if len(logArr) > numb {
			log = logArr[numb]
		}

		pass, _ := regexp.MatchString(line, log)

		fmt.Printf("%d: %t \n", numb+1, pass)
	}
}
