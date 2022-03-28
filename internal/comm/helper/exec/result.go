package execHelper

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	langHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/lang"
	"github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	websocketUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/websocket"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"github.com/mattn/go-runewidth"
	"regexp"
	"strconv"
	"strings"
)

func CheckCaseResult(scriptFile string, logs string, report *commDomain.ZtfReport, scriptIdx int,
	total int, secs string, pathMaxWidth int, numbMaxWidth int,
	wsMsg websocket.Message) {

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
		scriptIdx, total, secs, pathMaxWidth, numbMaxWidth, wsMsg)
}

func ValidateCaseResult(scriptFile string, langType string,
	steps []commDomain.ZentaoCaseStep, skip bool, actualArr [][]string, report *commDomain.ZtfReport,
	scriptIdx int, total int, secs string, pathMaxWidth int, numbMaxWidth int,
	wsMsg websocket.Message) {

	_, caseId, productId, title := scriptHelper.GetCaseInfo(scriptFile)

	stepLogs := make([]commDomain.StepLog, 0)
	caseResult := commConsts.PASS
	noExpects := true

	if skip {
		caseResult = commConsts.SKIP
	} else {
		stepIdxToCheck := 0
		for _, step := range steps { // iterate by checkpoints
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

	cs := commDomain.FuncResult{Id: caseId, ProductId: productId, Title: title,
		Path: scriptFile, Status: caseResult, Steps: stepLogs}
	report.FuncResult = append(report.FuncResult, cs)

	width := strconv.Itoa(len(strconv.Itoa(total)))
	numbWidth := strconv.Itoa(numbMaxWidth)

	path := cs.Path
	lent := runewidth.StringWidth(path)

	if pathMaxWidth > lent {
		postFix := strings.Repeat(" ", pathMaxWidth-lent)
		path += postFix
	}

	format := "(%" + width + "d/%d) %s [%s] [%" + numbWidth + "d. %s] (%ss)"

	status := i118Utils.Sprintf(cs.Status.String())
	msg := fmt.Sprintf(format, scriptIdx+1, total, status, path, cs.Id, cs.Title, secs)

	if commConsts.ComeFrom != "cmd" {
		websocketUtils.SendExecMsg(msg, "", wsMsg)
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
		if expect[:1] == "`" && expect[len(expect)-1:] == "`" {
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
