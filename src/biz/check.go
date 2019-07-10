package biz

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"utils"
)

func CheckResults(dir string, langType string) {
	scriptFiles, _ := utils.GetAllFiles(dir, langType)

	resultMap := make(map[string]bool)
	checkpointMap := make(map[string][]string)

	for _, scriptFile := range scriptFiles {
		logFile := utils.ScriptToLog(scriptFile)

		expectContent := utils.ReadExpect(scriptFile)
		logContent := utils.ReadFile(logFile)

		expectContent = strings.Trim(expectContent, "\n")
		logContent = strings.Trim(logContent, "\n")

		Compare(scriptFile, expectContent, logContent, &resultMap, &checkpointMap)
	}

	Print(resultMap, checkpointMap)
}

func Compare(scriptFile string, expectContent string, logContent string,
	resultMap *map[string]bool, checkpointMap *map[string][]string) {
	expectArr := strings.Split(expectContent, "\n")
	logArr := strings.Split(logContent, "\n")

	checkpoints := make([]string, 0)

	result := true

	for numb, line := range expectArr {
		log := "N/A"
		if len(logArr) > numb {
			log = logArr[numb]
		}

		pass, _ := regexp.MatchString(line, log)

		if !pass {
			result = false
		}

		checkpoints = append(checkpoints, "Line "+strconv.Itoa(numb+1)+": "+strconv.FormatBool(result))

		if !pass {
			checkpoints = append(checkpoints, "Expect "+line)
			checkpoints = append(checkpoints, "Actual "+log)
		}
	}

	(*resultMap)[scriptFile] = result
	(*checkpointMap)[scriptFile] = checkpoints
}

func Print(resultMap map[string]bool, checkpointMap map[string][]string) {
	fmt.Printf("")

	for script, result := range resultMap {

		fmt.Printf("\n=== Case %s: %t \n", script, result)
		if !result {
			checkpoints := checkpointMap[script]

			for _, line := range checkpoints {
				if strings.Index(line, "Line") > -1 {
					fmt.Printf("\n")
				}

				fmt.Printf("    %s \n", line)
			}
		}
	}
}
