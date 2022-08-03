package execHelper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	langHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/lang"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/mattn/go-runewidth"
)

func CheckCaseResult(scriptFile string, logs string, report *commDomain.ZtfReport, scriptIdx int,
	total int, secs string, pathMaxWidth int, numbMaxWidth int,
	wsMsg *websocket.Message, errOutput string) {

	steps, isOldFormat := scriptHelper.GetStepAndExpectMap(scriptFile)

	isIndependent, expectIndependentContent := scriptHelper.GetDependentExpect(scriptFile)
	if isIndependent {
		if isOldFormat {
			scriptHelper.GetExpectMapFromIndependentFileObsolete(&steps, expectIndependentContent, false)
		} else {
			scriptHelper.GetExpectMapFromIndependentFile(&steps, expectIndependentContent, false)
		}
	}

	skip := false
	actualArr := make([][]string, 0)
	if isOldFormat {
		skip, actualArr = scriptHelper.ReadLogArrObsolete(logs)
	} else {
		skip, actualArr = scriptHelper.ReadLogArr(logs)
	}

	language := langHelper.GetLangByFile(scriptFile)
	ValidateCaseResult(scriptFile, language, steps, skip, actualArr, report,
		scriptIdx, total, secs, pathMaxWidth, numbMaxWidth, wsMsg, errOutput)
}

func ValidateCaseResult(scriptFile string, langType string,
	steps []commDomain.ZentaoCaseStep, skip bool, actualArr [][]string, report *commDomain.ZtfReport,
	scriptIdx int, total int, secs string, pathMaxWidth int, numbMaxWidth int,
	wsMsg *websocket.Message, errOutput string) {

	key := stringUtils.Md5(scriptFile)

	_, caseId, productId, title := scriptHelper.GetCaseInfo(scriptFile)

	stepLogs := make([]commDomain.StepLog, 0)
	caseResult := commConsts.PASS
	noExpects := true

	if skip {
		caseResult = commConsts.SKIP
	} else {
		stepIdxToCheck := 0
		for index, step := range steps { // iterate by checkpoints
			stepName := strings.TrimSpace(step.Desc)
			expect := strings.TrimSpace(step.Expect)

			if expect == "" {
				continue
			}

			noExpects = false

			expectLines := strings.Split(expect, "\n")
			var actualLines []string
			if len(actualArr) > stepIdxToCheck {
				actualLines = actualArr[stepIdxToCheck]
			}

			stepResult, checkpointLogs := ValidateStepResult(langType, expectLines, actualLines)
			if errOutput != "" && index == 0 && len(checkpointLogs) > 0 {
				checkpointLogs[0].Actual = errOutput
			}
			stepLog := commDomain.StepLog{Id: strconv.Itoa(stepIdxToCheck), Name: stepName, Status: stepResult, CheckPoints: checkpointLogs}
			stepLogs = append(stepLogs, stepLog)
			if stepResult == commConsts.FAIL {
				caseResult = commConsts.FAIL
			}

			stepIdxToCheck++
		}
	}

	if noExpects {
		caseResult = commConsts.SKIP
	}

	if caseResult == commConsts.FAIL {
		report.Fail = report.Fail + 1
	} else if caseResult == commConsts.PASS {
		report.Pass = report.Pass + 1
	} else if caseResult == commConsts.SKIP {
		report.Skip = report.Skip + 1
	}
	report.Total = report.Total + 1

	relativePath := strings.TrimLeft(scriptFile, report.WorkspacePath)
	csResult := commDomain.FuncResult{Id: caseId, ProductId: productId, Title: title,
		Key: key, Path: scriptFile, RelativePath: relativePath, Status: caseResult, Steps: stepLogs}
	report.FuncResult = append(report.FuncResult, csResult)

	width := strconv.Itoa(len(strconv.Itoa(total)))
	numbWidth := strconv.Itoa(numbMaxWidth)

	path := csResult.Path
	lent := runewidth.StringWidth(path)

	if pathMaxWidth > lent {
		postFix := strings.Repeat(" ", pathMaxWidth-lent)
		path += postFix
	}

	format := "(%" + width + "d/%d) %s [%s] [%" + numbWidth + "d. %s] (%ss)"

	status := i118Utils.Sprintf(csResult.Status.String())
	msg := fmt.Sprintf(format, scriptIdx+1, total, status, path, csResult.Id, csResult.Title, secs)

	// print each case result
	if commConsts.ExecFrom != commConsts.FromCmd {
		msgCategory := commConsts.Result
		if csResult.Status == commConsts.FAIL {
			msgCategory = commConsts.Error
		}

		totalStepCount := len(csResult.Steps)
		passStepCount := 0
		failStepCount := 0

		failedCheckpoints := make([]string, 0)
		passStepCount, failStepCount = appendFailedStepResult(csResult, &failedCheckpoints)

		arr := []string{i118Utils.Sprintf("steps_result_msg", totalStepCount, passStepCount, failStepCount)}
		arr = append(arr, failedCheckpoints...)
		msg = strings.Join(arr, "<br/>")

		websocketHelper.SendExecMsg(msg, "", msgCategory,
			iris.Map{"key": key, "status": csResult.Status}, wsMsg)
	}
	logUtils.ExecConsole(color.FgCyan, msg)
	logUtils.ExecResult(msg)
}

func ValidateStepResult(langType string, expectLines []string, actualLines []string) (
	stepResult commConsts.ResultStatus, checkpointLogs []commDomain.CheckPointLog) {
	stepResult = commConsts.PASS

	idx := 0
	for _, expect := range expectLines {
		log := "N/A"
		if len(actualLines) > idx {
			log = strings.TrimSpace(actualLines[idx])
		}

		expect = strings.TrimSpace(expect)

		var pass bool
		if len(expect) >= 2 && expect[:1] == "`" && expect[len(expect)-1:] == "`" {
			expect = expect[1 : len(expect)-1]
			pass = MatchString(expect, log, langType)
		} else {
			pass = strings.Contains(log, expect)
		}

		if !pass {
			stepResult = commConsts.FAIL
		}

		cp := commDomain.CheckPointLog{Numb: idx + 1, Status: stepResult, Expect: expect, Actual: log}
		checkpointLogs = append(checkpointLogs, cp)

		idx++
	}

	return stepResult, checkpointLogs

}

func MatchString(expect string, actual string, langType string) bool {
	expect = strings.TrimSpace(expect)
	actual = strings.TrimSpace(actual)

	expect = strings.Replace(expect, "%s", `.+?`, -1)                                  // 字符串
	expect = strings.Replace(expect, "%i", `[+\-]?[0-9]+`, -1)                         // 十进制数字，可有符号
	expect = strings.Replace(expect, "%d", `[0-9]+`, -1)                               // 十进制数字，无符号
	expect = strings.Replace(expect, "%x", `(0X|0x)?[0-9a-fA-F]+`, -1)                 // 十六进制数字
	expect = strings.Replace(expect, "%f", `[+\-]?\.?[0-9]+\.?[0-9]*(E-?[0-9]+)?`, -1) // 十进制浮点数
	expect = strings.Replace(expect, "%c", ".", -1)                                    // 单个字符

	pass, _ := regexp.MatchString(expect, actual)
	return pass
}
