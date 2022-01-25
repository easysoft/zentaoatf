package scriptUtils

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/script"
	"github.com/emirpasic/gods/maps"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"github.com/mattn/go-runewidth"
	"regexp"
	"strconv"
	"strings"
)

func CheckCaseResult(scriptFile string, logs string, report *commDomain.ZtfReport, idx int,
	total int, secs string, pathMaxWidth int, numbMaxWidth int,
	sendOutMsg, sendExecMsg func(info, isRunning string, wsMsg websocket.Message), wsMsg websocket.Message) {

	stepMap, _, expectMap, isOldFormat := scriptUtils.GetStepAndExpectMap(scriptFile)

	isIndependent, expectIndependentContent := scriptUtils.GetDependentExpect(scriptFile)
	if isIndependent {
		if isOldFormat {
			expectMap = scriptUtils.GetExpectMapFromIndependentFileObsolete(expectMap, expectIndependentContent, false)
		} else {
			expectMap = scriptUtils.GetExpectMapFromIndependentFile(expectMap, expectIndependentContent, false)
		}
	}

	skip := false
	actualArr := make([][]string, 0)
	if isOldFormat {
		skip, actualArr = scriptUtils.ReadLogArrObsolete(logs)
	} else {
		skip, actualArr = scriptUtils.ReadLogArr(logs)
	}

	language := langUtils.GetLangByFile(scriptFile)
	ValidateCaseResult(scriptFile, language, stepMap, expectMap, skip, actualArr, report,
		idx, total, secs, pathMaxWidth, numbMaxWidth,
		sendOutMsg, sendExecMsg, wsMsg)
}

func ValidateCaseResult(scriptFile string, langType string,
	stepMap, expectMap maps.Map, skip bool, actualArr [][]string, report *commDomain.ZtfReport,
	idx int, total int, secs string, pathMaxWidth int, numbMaxWidth int,
	sendOutMsg, sendExecMsg func(info, isRunning string, wsMsg websocket.Message), wsMsg websocket.Message) {

	_, caseId, productId, title := scriptUtils.GetCaseInfo(scriptFile)

	stepLogs := make([]commDomain.StepLog, 0)
	caseResult := commConsts.PASS
	noExpects := true

	if skip {
		caseResult = commConsts.SKIP
	} else {
		idx := 0

		for _, numbObj := range expectMap.Keys() { // iterate by checkpoints
			stepObj, _ := stepMap.Get(numbObj)
			expectObj, _ := expectMap.Get(numbObj)

			numb := strings.TrimSpace(numbObj.(string))
			stepName := strings.TrimSpace(stepObj.(string))
			expect := strings.TrimSpace(expectObj.(string))

			if expect == "" {
				continue
			}

			noExpects = false

			expectLines := strings.Split(expect, "\n")
			var actualLines []string
			if len(actualArr) > idx {
				actualLines = actualArr[idx]
			}

			stepResult, checkpointLogs := ValidateStepResult(langType, expectLines, actualLines)
			stepLog := commDomain.StepLog{Id: strings.TrimRight(numb, "."), Name: stepName, Status: stepResult, CheckPoints: checkpointLogs}
			stepLogs = append(stepLogs, stepLog)
			if stepResult == commConsts.FAIL {
				caseResult = commConsts.FAIL
			}

			idx++
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
	msg := fmt.Sprintf(format, idx+1, total, status, path, cs.Id, cs.Title, secs)

	sendExecMsg(msg, "", wsMsg)
	logUtils.ExecConsole(color.FgCyan, msg)
	logUtils.ExecResult(msg)
}

func ValidateStepResult(langType string, expectLines []string, actualLines []string) (
	stepResult commConsts.ResultStatus, checkpointLogs []commDomain.CheckPointLog) {
	stepResult = commConsts.PASS

	indx2 := 0
	for _, expect := range expectLines {
		log := "N/A"
		if len(actualLines) > indx2 {
			log = actualLines[indx2]
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

		cp := commDomain.CheckPointLog{Numb: indx2 + 1, Status: stepResult, Expect: expect, Actual: log}
		checkpointLogs = append(checkpointLogs, cp)

		indx2++
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
