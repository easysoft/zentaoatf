package biz

import (
	"fmt"
	"model"
	"regexp"
	"strings"
	"time"
	"utils"
)

func CheckResults(dir string, langType string, report *model.TestReport) {
	fmt.Printf("\n=== Begin to analyse test result ===\n\n")

	scriptFiles, _ := utils.GetAllFiles(dir, langType)

	for _, scriptFile := range scriptFiles {
		logFile := utils.ScriptToLogName(dir, scriptFile)

		expectContent := utils.ReadExpect(scriptFile)
		logContent := utils.ReadFile(logFile)

		expectContent = strings.Trim(expectContent, "\n")
		logContent = strings.Trim(logContent, "\n")

		Compare(scriptFile, expectContent, logContent, report)
	}
}

func Compare(scriptFile string, expectContent string, logContent string, report *model.TestReport) {
	expectArr := strings.Split(expectContent, "\n")
	logArr := strings.Split(logContent, "\n")

	checkpoints := make([]model.CheckPointLog, 0)

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

		cp := model.CheckPointLog{Numb: numb + 1, Status: result, Expect: line, Actual: log}
		checkpoints = append(checkpoints, cp)
	}

	if !result {
		report.Fail = report.Fail + 1
	} else {
		report.Pass = report.Pass + 1
	}
	report.Total = report.Total + 1

	cs := model.CaseLog{Path: scriptFile, Status: result, CheckPoints: checkpoints}
	report.Cases = append(report.Cases, cs)
}

func Print(report model.TestReport, workDir string) {
	startSec := time.Unix(report.StartTime, 0)
	endSec := time.Unix(report.EndTime, 0)

	logs := make([]string, 0)

	PrintAndLog(&logs, fmt.Sprintf("Run scripts in folder \"%s\" on %s OS\n",
		report.Path, report.Env))

	PrintAndLog(&logs, fmt.Sprintf("From %s to %s, duration %d sec",
		startSec.Format("2006-01-02 15:04:05"), endSec.Format("2006-01-02 15:04:05"), report.Duration))

	PrintAndLog(&logs, fmt.Sprintf("Total: %d, Fail: %d, Pass: %d",
		report.Total, report.Pass, report.Fail))

	for _, cs := range report.Cases {
		PrintAndLog(&logs, fmt.Sprintf("\n%s: %t", cs.Path, cs.Status))

		if len(cs.CheckPoints) > 0 {
			count := 0
			for _, cp := range cs.CheckPoints {
				if count > 0 {
					PrintAndLog(&logs, "")
				}

				PrintAndLog(&logs, fmt.Sprintf("   Line %d: %t", cp.Numb, cp.Status))
				PrintAndLog(&logs, fmt.Sprintf("   Expect %s", cp.Expect))
				PrintAndLog(&logs, fmt.Sprintf("   Actual %s", cp.Actual))

				count++
			}
		} else {
			PrintAndLog(&logs, "   No checkpoints")
		}
	}

	utils.WriteFile(workDir+"/logs/report.txt", strings.Join(logs, "\n"))
}

func PrintAndLog(logs *[]string, str string) {
	*logs = append(*logs, str)
	fmt.Println(str)
}
