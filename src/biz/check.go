package biz

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"utils"
)

func CheckResults(dir string, langType string,
	summaryMap *map[string]interface{}, resultMap *map[string]bool, checkpointMap *map[string][]string) {
	fmt.Printf("\n=== Begin to analyse test result ===\n\n")

	(*summaryMap)["pass"] = 0
	(*summaryMap)["fail"] = 0
	(*summaryMap)["total"] = 0

	scriptFiles, _ := utils.GetAllFiles(dir, langType)

	for _, scriptFile := range scriptFiles {
		logFile := utils.ScriptToLog(scriptFile)

		expectContent := utils.ReadExpect(scriptFile)
		logContent := utils.ReadFile(logFile)

		expectContent = strings.Trim(expectContent, "\n")
		logContent = strings.Trim(logContent, "\n")

		Compare(scriptFile, expectContent, logContent, summaryMap, resultMap, checkpointMap)
	}
}

func Compare(scriptFile string, expectContent string, logContent string,
	summaryMap *map[string]interface{}, resultMap *map[string]bool, checkpointMap *map[string][]string) {
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
			(*summaryMap)["fail"] = (*summaryMap)["fail"].(int) + 1
		} else {
			(*summaryMap)["pass"] = (*summaryMap)["pass"].(int) + 1
		}
		(*summaryMap)["total"] = (*summaryMap)["total"].(int) + 1

		checkpoints = append(checkpoints, "Line "+strconv.Itoa(numb+1)+": "+strconv.FormatBool(result))

		if !pass {
			checkpoints = append(checkpoints, "Expect "+line)
			checkpoints = append(checkpoints, "Actual "+log)
		}
	}

	(*resultMap)[scriptFile] = result
	(*checkpointMap)[scriptFile] = checkpoints
}

func Print(summaryMap map[string]interface{}, resultMap map[string]bool, checkpointMap map[string][]string) {
	startSec := time.Unix(summaryMap["startTime"].(int64), 0)
	endSec := time.Unix(summaryMap["endTime"].(int64), 0)

	fmt.Printf("From %s to %s, duration %d sec \n",
		startSec.Format("2006-01-02 15:04:05"),
		endSec.Format("2006-01-02 15:04:05"),
		summaryMap["duration"])

	fmt.Printf("Total: %d, Fail: %d, Pass: %d \n",
		summaryMap["total"], summaryMap["pass"], summaryMap["fail"])

	for script, result := range resultMap {

		fmt.Printf("\n--- Case %s: %t \n", script, result)
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
