package execHelper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

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

func CheckCaseResult(scriptFile string, logs string, report *commDomain.ZtfReport,
	total int, secs string, pathMaxWidth int, numbMaxWidth int,
	wsMsg *websocket.Message, errOutput string, lock *sync.Mutex) {

	steps := scriptHelper.GetStepAndExpectMap(scriptFile)

	isIndependent, expectIndependentContent := scriptHelper.GetDependentExpect(scriptFile)
	if isIndependent {
		scriptHelper.GetExpectMapFromIndependentFile(&steps, expectIndependentContent, false)
	}

	skip := false
	actualArr := make([][]string, 0)
	skip, actualArr = scriptHelper.ReadLogArr(logs)
	if len(actualArr) == 0 {
		skip, actualArr = scriptHelper.ReadLogArrOld(logs)
	}

	language := langHelper.GetLangByFile(scriptFile)
	ValidateCaseResult(scriptFile, language, steps, skip, actualArr, report,
		total, secs, pathMaxWidth, numbMaxWidth, wsMsg, errOutput, lock)
}

func ValidateCaseResult(scriptFile string, langType string,
	steps []commDomain.ZentaoCaseStep, skip bool, actualArr [][]string, report *commDomain.ZtfReport, total int, secs string, pathMaxWidth int, numbMaxWidth int,
	wsMsg *websocket.Message, errOutput string, lock *sync.Mutex) {

	key := stringUtils.Md5(scriptFile)

	_, caseId, productId, title, _ := scriptHelper.GetCaseInfo(scriptFile)

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
			stepLog := commDomain.StepLog{Id: strconv.Itoa(stepIdxToCheck + 1), Name: stepName, Status: stepResult, CheckPoints: checkpointLogs}
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

	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
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
	msg := fmt.Sprintf(format, len(report.FuncResult), total, status, path, csResult.Id, csResult.Title, secs)

	// print each case result
	if commConsts.ExecFrom == commConsts.FromClient {
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
		if len(expect) > 1 && expect[0:1] == "~" && expect[len(expect)-1:] == "~" {
			pass = MatchScene(expect[1:len(expect)-1], log, langType)
		} else if len(expect) >= 2 && expect[:1] == "`" && expect[len(expect)-1:] == "`" {
			expect = expect[1 : len(expect)-1]
			pass = MatchString(expect, log, langType)
		} else {
			pass = strings.TrimSpace(log) == strings.TrimSpace(expect)
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

func MatchScene(expect, actual, langType string) (pass bool) {
	expect = strings.TrimSpace(expect)
	actual = strings.TrimSpace(actual)
	if len(expect) == 0 {
		return actual == ""
	}

	if len(expect) > 2 {
		scene := expect[:2]
		expect = strings.TrimSpace(expect[2:])
		switch scene {
		case "f:":
			if strings.Contains(expect, "*") {
				expectArr := strings.Split(expect, "*")
				repeatCount, _ := strconv.Atoi(string(expectArr[1]))
				return strings.Count(actual, expectArr[0]) >= repeatCount
			}
			if expect[0:1] == "(" && expect[len(expect)-1:] == ")" && strings.Contains(expect, ",") {
				expect = fmt.Sprintf("^%s{1}$", strings.ReplaceAll(expect, ",", "|"))
			}
			return MatchString(expect, actual, langType)
		case "m:":
			return MatchString(expect, actual, langType)
		case "c:":
			if len(expect) > 2 && (expect[:2] == ">=" || expect[:2] == "<=" || expect[:2] == "<>" || expect[:2] == "!=") {
				character := expect[:2]
				expectFloot, err := strconv.ParseFloat(strings.TrimSpace(expect[2:]), 64)
				if err != nil {
					return false
				}
				actualFloot, err := strconv.ParseFloat(strings.TrimSpace(actual), 64)
				if err != nil {
					return false
				}
				switch character {
				case ">=":
					return actualFloot >= expectFloot
				case "<=":
					return actualFloot <= expectFloot
				case "<>":
					return actualFloot != expectFloot
				case "!=":
					return actualFloot != expectFloot
				}
			} else if strings.Contains(expect, "-") && strings.Count(expect, "-") == 1 {
				rangeArr := strings.Split(expect, "-")
				rangeFrom, err := strconv.ParseFloat(strings.TrimSpace(rangeArr[0]), 64)
				if err != nil {
					return false
				}
				rangeTo, err := strconv.ParseFloat(strings.TrimSpace(rangeArr[1]), 64)
				if err != nil {
					return false
				}
				actualFloot, err := strconv.ParseFloat(strings.TrimSpace(actual), 64)
				if err != nil {
					return false
				}
				return actualFloot >= rangeFrom && actualFloot <= rangeTo
			} else {
				character := expect[:1]
				expectFloot, err := strconv.ParseFloat(strings.TrimSpace(expect[1:]), 64)
				if err != nil {
					return false
				}
				actualFloot, err := strconv.ParseFloat(strings.TrimSpace(actual), 64)
				if err != nil {
					return false
				}
				switch character {
				case ">":
					return actualFloot > expectFloot
				case "<":
					return actualFloot < expectFloot
				case "=":
					return actualFloot == expectFloot
				}
				if strings.Contains(expect, "-") {
					expectArr := strings.Split(expect, "-")
					expectMin, _ := strconv.ParseFloat(strings.TrimSpace(expectArr[0]), 64)
					expectMax, _ := strconv.ParseFloat(strings.TrimSpace(expectArr[1]), 64)
					return actualFloot >= expectMin && actualFloot <= expectMax
				}
			}
		case "l:":
			openParenthesisCount, closeParenthesisCount := 0, 0
			hasLogicCharacter := false
			for index, v := range expect {
				if v == '(' {
					openParenthesisCount++
				} else if v == ')' {
					closeParenthesisCount++
				}
				if v == '&' || v == '|' {
					hasLogicCharacter = true
				}
				if v == '|' && index > 0 && expect[index-1] != '\\' && (openParenthesisCount == closeParenthesisCount) {
					return MatchScene("l:"+expect[:index], actual, langType) || MatchScene("l:"+expect[index+1:], actual, langType)
				} else if v == '&' && index > 0 && string(expect[index-1]) != "\\" && (openParenthesisCount == closeParenthesisCount) {
					return MatchScene("l:"+expect[:index], actual, langType) && MatchScene("l:"+expect[index+1:], actual, langType)
				}
			}
			if expect[:1] == "(" {
				expect = expect[1:]
			}
			if expect[len(expect)-1:] == ")" {
				expect = expect[:len(expect)-1]
			}
			if !hasLogicCharacter {
				if expect[:1] == "!" {
					return !MatchScene(expect[1:], actual, langType)
				}
				return MatchScene(expect, actual, langType)
			} else {
				return MatchScene("l:"+expect, actual, langType)
			}
		}
	}

	return pass
}
