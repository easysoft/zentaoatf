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
		logFile := utils.ScriptToLogName(scriptFile)

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
		line = strings.TrimSpace(line)
		if line == "#" || line == "" {
			continue
		}

		log := "N/A"
		if len(logArr) > numb {
			log = logArr[numb]
			log = strings.TrimSpace(log)
		}

		pass, _ := regexp.MatchString(line, log)

		if !pass {
			result = false
		}

		checkpoints = append(checkpoints, "Line "+strconv.Itoa(numb+1)+": "+strconv.FormatBool(result))

		checkpoints = append(checkpoints, "Expect "+line)
		checkpoints = append(checkpoints, "Actual "+log)
	}

	if !result {
		(*summaryMap)["fail"] = (*summaryMap)["fail"].(int) + 1
	} else {
		(*summaryMap)["pass"] = (*summaryMap)["pass"].(int) + 1
	}
	(*summaryMap)["total"] = (*summaryMap)["total"].(int) + 1

	(*resultMap)[scriptFile] = result
	(*checkpointMap)[scriptFile] = checkpoints
}

func Print(summaryMap map[string]interface{}, resultMap map[string]bool, checkpointMap map[string][]string, workDir string) {
	startSec := time.Unix(summaryMap["startTime"].(int64), 0)
	endSec := time.Unix(summaryMap["endTime"].(int64), 0)

	var log string
	logs := make([]string, 0)

	log = fmt.Sprintf("From %s to %s, duration %d sec",
		startSec.Format("2006-01-02 15:04:05"),
		endSec.Format("2006-01-02 15:04:05"),
		summaryMap["duration"])
	logs = append(logs, log)
	fmt.Println(log)

	log = fmt.Sprintf("Total: %d, Fail: %d, Pass: %d",
		summaryMap["total"], summaryMap["pass"], summaryMap["fail"])
	logs = append(logs, log)
	fmt.Println(log)

	for script, result := range resultMap {
		count := 0
		log = fmt.Sprintf("\n--- Case %s: %t", script, result)
		logs = append(logs, log)
		fmt.Println(log)

		checkpoints := checkpointMap[script]

		for _, line := range checkpoints {
			if count > 0 && strings.Index(line, "Line ") > -1 {
				logs = append(logs, "\n")
				fmt.Println("")
			}

			log = fmt.Sprintf("    %s", line)
			logs = append(logs, log)
			fmt.Println(log)

			count++
		}
	}

	utils.WriteFile(workDir+"/logs/log.txt", strings.Join(logs, "\n"))
}
