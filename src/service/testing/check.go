package testingService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils/const"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/string"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func CheckCaseResult(scriptFile string, report *model.TestReport, idx int, total int, secs string, pathMaxWidth int) {
	logFile := zentaoUtils.ScriptToLogName(scriptFile)

	checkpointStepArr := zentaoUtils.ReadCheckpointSteps(scriptFile)
	expectArr := zentaoUtils.ReadExpect(scriptFile)
	skip, logArr := zentaoUtils.ReadLog(logFile)

	language := ""
	ValidateCaseResult(scriptFile, language, checkpointStepArr, expectArr, skip, logArr, report, idx, total, secs, pathMaxWidth)
}

func ValidateCaseResult(scriptFile string, langType string,
	checkpointStepArr []string, expectArr [][]string, skip bool, actualArr [][]string, report *model.TestReport,
	idx int, total int, secs string, pathMaxWidth int) {

	_, caseId, productId, title := zentaoUtils.GetCaseInfo(scriptFile)

	stepLogs := make([]model.StepLog, 0)
	caseResult := constant.PASS.String()

	if skip {
		caseResult = constant.SKIP.String()
	} else {
		indx := 0
		for _, step := range checkpointStepArr { // iterate by checkpoints
			var expectLines []string
			var actualLines []string

			if len(expectArr) > indx {
				expectLines = expectArr[indx]
			}
			if len(actualArr) > indx {
				actualLines = actualArr[indx]
			}

			re, _ := regexp.Compile(`\s{2,}`)
			step = re.ReplaceAllString(step, " ") // 多个空格替换成一个

			arr := strings.Split(step, " ")
			str := strings.Replace(arr[0], "@step", "", -1)
			stepId, _ := strconv.Atoi(str)

			stepResult, checkpointLogs := ValidateStepResult(langType, expectLines, actualLines)
			stepLog := model.StepLog{Id: stepId, Name: arr[1], Status: stepResult, CheckPoints: checkpointLogs}
			stepLogs = append(stepLogs, stepLog)
			if !stepResult {
				caseResult = constant.FAIL.String()
			}

			indx++
		}
	}

	if caseResult == constant.FAIL.String() {
		report.Fail = report.Fail + 1
	} else if caseResult == constant.PASS.String() {
		report.Pass = report.Pass + 1
	} else if caseResult == constant.SKIP.String() {
		report.Skip = report.Skip + 1
	}
	report.Total = report.Total + 1

	cs := model.CaseLog{Id: caseId, ProductId: productId, Title: title,
		Path: scriptFile, Status: caseResult, Steps: stepLogs}
	report.Cases = append(report.Cases, cs)

	// print case result to console
	statusColor := logUtils.ColoredStatus(cs.Status)
	width := strconv.Itoa(len(strconv.Itoa(total)))

	path := cs.Path
	lent := utf8.RuneCountInString(path)

	if pathMaxWidth > lent {
		postFix := strings.Repeat(" ", pathMaxWidth-lent)
		path += postFix
	}

	logUtils.Screen(
		fmt.Sprintf("(%"+width+"d/%d) %s [%s] %d %s (%ss)", idx+1, total, statusColor, path, cs.Id, cs.Title, secs))
}

func ValidateStepResult(langType string, expectLines []string, actualLines []string) (bool, []model.CheckPointLog) {
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
