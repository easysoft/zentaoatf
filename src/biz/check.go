package biz

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
	"strings"
	"time"
)

func CheckResults(dir string, langType string, report *model.TestReport) {
	fmt.Printf("\n=== Begin to analyse test result ===\n\n")

	scriptFiles, _ := utils.GetAllFiles(dir, langType)

	for _, scriptFile := range scriptFiles {
		logFile := utils.ScriptToLogName(dir, scriptFile)

		stepArr := utils.ReadCheckpointSteps(scriptFile)
		expectArr := utils.ReadExpect(scriptFile)
		logArr := utils.ReadLog(logFile)

		ValidateTestCase(scriptFile, langType, stepArr, expectArr, logArr, report)
	}
}

func ValidateTestCase(scriptFile string, langType string,
	stepArr []string, expectArr [][]string, actualArr [][]string, report *model.TestReport) {

	stepLogs := make([]model.StepLog, 0)

	caseResult := true

	indx := 0
	for _, step := range stepArr {
		stepResult := true

		var expectLines []string
		if len(expectArr) > indx {
			expectLines = expectArr[indx]
		}

		var actualLines []string
		if len(actualArr) > indx {
			actualLines = actualArr[indx]
		}

		checkpointLogs := make([]model.CheckPointLog, 0)

		indx2 := 0
		for _, expect := range expectLines {
			log := "N/A"
			if len(actualLines) > indx2 {
				log = actualLines[indx2]
			}

			pass := MatchString(expect, log, langType)
			if !pass {
				stepResult = false
			}

			cp := model.CheckPointLog{Numb: indx2 + 1, Status: pass, Expect: expect, Actual: log}
			checkpointLogs = append(checkpointLogs, cp)

			indx2++
		}

		step := model.StepLog{Numb: indx + 1, Name: step, Status: stepResult, CheckPoints: checkpointLogs}
		stepLogs = append(stepLogs, step)

		indx++
	}

	if !caseResult {
		report.Fail = report.Fail + 1
	} else {
		report.Pass = report.Pass + 1
	}
	report.Total = report.Total + 1

	cs := model.CaseLog{Path: scriptFile, Status: caseResult, Steps: stepLogs}
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

		if len(cs.Steps) > 0 {
			count := 0
			for _, step := range cs.Steps {
				if count > 0 { // 空行
					PrintAndLog(&logs, "")
				}

				PrintAndLog(&logs, fmt.Sprintf("  Step %d %s: %t", step.Numb, step.Name, step.Status))

				count1 := 0
				for _, cp := range step.CheckPoints {
					if count1 > 0 { // 空行
						PrintAndLog(&logs, "")
					}

					PrintAndLog(&logs, fmt.Sprintf("    Checkpoint %d: %t", cp.Numb, cp.Status))
					PrintAndLog(&logs, fmt.Sprintf("      Expect %s", cp.Expect))
					PrintAndLog(&logs, fmt.Sprintf("      Actual %s", cp.Actual))

					count1++
				}

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
