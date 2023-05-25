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
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/mattn/go-runewidth"
)

func CheckCaseResult(execParams commDomain.ExecParams, logs string, wsMsg *websocket.Message, errOutput string, lock *sync.Mutex) {
	steps := scriptHelper.GetStepAndExpectMap(execParams.ScriptFile)

	isIndependent, expectIndependentContent := scriptHelper.GetDependentExpect(execParams.ScriptFile)
	if isIndependent {
		scriptHelper.GetExpectMapFromIndependentFile(&steps, expectIndependentContent, false)
	}

	skip := false
	skip, actualArr := scriptHelper.ReadLogArr(logs)
	if len(actualArr) == 0 {
		skip, actualArr = scriptHelper.ReadLogArrOld(logs)
	}

	language := langHelper.GetLangByFile(execParams.ScriptFile)
	ValidateCaseResult(execParams, language, steps, skip, actualArr, wsMsg, errOutput, lock)
}

func ValidateCaseResult(execParams commDomain.ExecParams, langType string,
	steps []commDomain.ZentaoCaseStep, skip bool, actualArr [][]string,
	wsMsg *websocket.Message, errOutput string, lock *sync.Mutex) {

	key := stringUtils.Md5(execParams.ScriptFile)

	_, caseId, productId, title, _ := scriptHelper.GetCaseInfo(execParams.ScriptFile)

	stepLogs, caseResult := getStepLogs(skip, steps, actualArr, langType, errOutput)

	if lock != nil {
		lock.Lock()
	}

	incrReportNum(caseResult, execParams.Report)

	relativePath := strings.TrimPrefix(execParams.ScriptFile, commConsts.WorkDir)
	csResult := commDomain.FuncResult{Id: caseId, ProductId: productId, Title: title,
		Key: key, Path: execParams.ScriptFile, RelativePath: relativePath, Status: caseResult, Steps: stepLogs}
	execParams.Report.FuncResult = append(execParams.Report.FuncResult, csResult)

	if lock != nil {
		lock.Unlock()
	}

	width := strconv.Itoa(len(strconv.Itoa(len(execParams.CasesToRun))))

	path := relativePath
	csTitle := csResult.Title
	lenp := runewidth.StringWidth(csResult.Path)
	lent := runewidth.StringWidth(csTitle)

	if execParams.PathMaxWidth > lenp {
		postFix := strings.Repeat(" ", execParams.PathMaxWidth-lenp)
		path += postFix
		relativePath += postFix
	}

	if execParams.TitleMaxWidth > lent {
		postFix := strings.Repeat(" ", execParams.TitleMaxWidth-lent)
		csTitle += postFix
	}

	format := "(%" + width + "d/%d) [%s] [%s] [%s] [%ss]"

	statusWithColor, status := GenStatusTxt(csResult.Status)

	msg := fmt.Sprintf(format, execParams.ScriptIdx+1, len(execParams.CasesToRun), status, path, csTitle, execParams.Secs)
	msgWithColor := fmt.Sprintf(format, execParams.ScriptIdx+1, len(execParams.CasesToRun), statusWithColor, path, csTitle, execParams.Secs)

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
	logUtils.ExecConsole(-1, msgWithColor)
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
			pass = Match(expect, log)
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

func MatchScene(expect, actual, langType string) (pass bool) {
	expect = strings.TrimSpace(expect)
	actual = strings.TrimSpace(actual)

	if len(expect) == 0 {
		pass = actual == ""
		return
	}

	if len(expect) <= 2 {
		return
	}

	// len(expect) > 2
	scene := expect[:2]
	expect = strings.TrimSpace(expect[2:])

	switch scene {
	case "f:":
		return Contain(expect, actual, langType)

	case "m:":
		return Match(expect, actual)

	case "c:":
		return Compare(expect, actual)

	case "l:":
		return Logic(expect, actual, langType)
	}

	return
}

func Contain(expect, actual string, langType string) bool {
	if strings.Contains(expect, "*") {
		expectArr := strings.Split(expect, "*")
		repeatCount, _ := strconv.Atoi(string(expectArr[1]))
		return strings.Count(actual, expectArr[0]) >= repeatCount
	}
	if expect[0:1] == "(" && expect[len(expect)-1:] == ")" && strings.Contains(expect, ",") {
		expect = fmt.Sprintf("^%s{1}$", strings.ReplaceAll(expect, ",", "|"))
	}

	return Match(expect, actual)
}

func Match(expect string, actual string) bool {
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

func Compare(expect string, actual string) bool {
	if len(expect) > 2 && stringUtils.FindInArr(expect[:2], []string{">=", "<=", "<>", "!="}) {
		character := expect[:2]
		expectFloat, err := strconv.ParseFloat(strings.TrimSpace(expect[2:]), 64)
		if err != nil {
			return false
		}
		actualFloat, err := strconv.ParseFloat(strings.TrimSpace(actual), 64)
		if err != nil {
			return false
		}

		compareResult := CompareFloat(actualFloat, expectFloat, character)

		return compareResult
	}

	if strings.Contains(expect, "-") && strings.Count(expect, "-") == 1 {
		return CompareRange(expect, actual)
	}

	character := expect[:1]
	expectFloat, err := strconv.ParseFloat(strings.TrimSpace(expect[1:]), 64)
	if err != nil {
		return false
	}
	actualFloat, err := strconv.ParseFloat(strings.TrimSpace(actual), 64)
	if err != nil {
		return false
	}

	switch character {
	case ">":
		return actualFloat > expectFloat
	case "<":
		return actualFloat < expectFloat
	case "=":
		return actualFloat == expectFloat
	}

	if strings.Contains(expect, "-") {
		expectArr := strings.Split(expect, "-")
		expectMin, _ := strconv.ParseFloat(strings.TrimSpace(expectArr[0]), 64)
		expectMax, _ := strconv.ParseFloat(strings.TrimSpace(expectArr[1]), 64)
		return actualFloat >= expectMin && actualFloat <= expectMax
	}

	return false
}

func CompareFloat(actualFloat, expectFloat float64, symbol string) bool {
	switch symbol {
	case ">=":
		return actualFloat >= expectFloat
	case "<=":
		return actualFloat <= expectFloat
	case "<>":
		return actualFloat != expectFloat
	case "!=":
		return actualFloat != expectFloat
	case ">":
		return actualFloat > expectFloat
	case "<":
		return actualFloat < expectFloat
	case "=":
		return actualFloat == expectFloat
	}

	return false
}

func CompareRange(expect string, actual string) bool {
	rangeArr := strings.Split(expect, "-")
	rangeFrom, err := strconv.ParseFloat(strings.TrimSpace(rangeArr[0]), 64)
	if err != nil {
		return false
	}

	rangeTo, err := strconv.ParseFloat(strings.TrimSpace(rangeArr[1]), 64)
	if err != nil {
		return false
	}

	actualFloat, err := strconv.ParseFloat(strings.TrimSpace(actual), 64)
	if err != nil {
		return false
	}

	return actualFloat >= rangeFrom && actualFloat <= rangeTo
}

func Logic(expect, actual, langType string) bool {
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

func getStepLogs(skip bool, steps []commDomain.ZentaoCaseStep, actualArr [][]string, langType string, errOutput string) (stepLogs []commDomain.StepLog, caseResult commConsts.ResultStatus) {
	stepLogs = make([]commDomain.StepLog, 0)
	caseResult = commConsts.PASS
	noExpects := true

	if skip {
		caseResult = commConsts.SKIP
		return
	}

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

	if noExpects {
		caseResult = commConsts.SKIP
	}

	return
}

func incrReportNum(caseResult commConsts.ResultStatus, report *commDomain.ZtfReport) {
	if caseResult == commConsts.FAIL {
		report.Fail = report.Fail + 1
	} else if caseResult == commConsts.PASS {
		report.Pass = report.Pass + 1
	} else if caseResult == commConsts.SKIP {
		report.Skip = report.Skip + 1
	}

	report.Total = report.Total + 1
}
