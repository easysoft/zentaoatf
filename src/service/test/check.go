package testingService

import (
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/string"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/fatih/color"
	"regexp"
)

func CheckResults(files []string, dir string, langType string, report *model.TestReport) {
	logUtils.Printt("\n")
	logUtils.PrintWholeLine(i118Utils.I118Prt.Sprintf("begin_analyse"), "=", color.FgCyan)

	for _, scriptFile := range files {
		logFile := zentaoUtils.ScriptToLogName(dir, scriptFile)

		stepArr := zentaoUtils.ReadCheckpointSteps(scriptFile)
		expectArr := zentaoUtils.ReadExpect(scriptFile)
		skip, logArr := zentaoUtils.ReadLog(logFile)

		ValidateTestCase(scriptFile, langType, stepArr, expectArr, skip, logArr, report)
	}
}

func ValidateTestCase(scriptFile string, langType string,
	stepArr []string, expectArr [][]string, skip bool, actualArr [][]string, report *model.TestReport) {

	stepLogs := make([]model.StepLog, 0)
	caseResult := constant.PASS

	if skip {
		caseResult = constant.SKIP
	} else {
		indx := 0
		for _, step := range stepArr {
			var expectLines []string
			if len(expectArr) > indx {
				expectLines = expectArr[indx]
			}

			var actualLines []string
			if len(actualArr) > indx {
				actualLines = actualArr[indx]
			}

			re, _ := regexp.Compile(`\s{2,}`)
			step = re.ReplaceAllString(step, " ")

			stepResult, checkpointLogs := ValidateStep(langType, expectLines, actualLines)
			stepLog := model.StepLog{Numb: indx + 1, Name: step, Status: stepResult, CheckPoints: checkpointLogs}
			stepLogs = append(stepLogs, stepLog)
			if !stepResult {
				caseResult = constant.FAIL
			}

			indx++
		}
	}

	if caseResult == constant.FAIL {
		report.Fail = report.Fail + 1
	} else if caseResult == constant.PASS {
		report.Pass = report.Pass + 1
	} else if caseResult == constant.SKIP {
		report.Skip = report.Skip + 1
	}
	report.Total = report.Total + 1

	cs := model.CaseLog{Path: scriptFile, Status: caseResult, Steps: stepLogs}
	report.Cases = append(report.Cases, cs)
}

func ValidateStep(langType string, expectLines []string, actualLines []string) (bool, []model.CheckPointLog) {
	stepResult := true

	checkpointLogs := make([]model.CheckPointLog, 0)

	indx2 := 0
	for _, expect := range expectLines {
		log := "N/A"
		if len(actualLines) > indx2 {
			log = actualLines[indx2]
		}

		pass := stringUtils.MatchString(expect, log, langType)
		if !pass {
			stepResult = false
		}

		cp := model.CheckPointLog{Numb: indx2 + 1, Status: pass, Expect: expect, Actual: log}
		checkpointLogs = append(checkpointLogs, cp)

		indx2++
	}

	return stepResult, checkpointLogs

}
