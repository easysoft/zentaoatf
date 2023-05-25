package execHelper

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	dateUtils "github.com/easysoft/zentaoatf/pkg/lib/date"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	shellUtils "github.com/easysoft/zentaoatf/pkg/lib/shell"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/mattn/go-runewidth"
)

func GenUnitTestReport(req serverDomain.TestSet, startTime, endTime int64, ch chan int, wsMsg *websocket.Message) (
	report commDomain.ZtfReport) {

	key := stringUtils.Md5(req.WorkspacePath)

	testSuites := RetrieveUnitResult(req, startTime)
	ZipDir(req)

	cases, classNameMaxWidth, duration := ParserUnitTestResult(testSuites)

	if duration == 0 {
		duration = float32(endTime - startTime)
	}

	report = commDomain.ZtfReport{
		Name:        req.Name,
		Platform:    commonUtils.GetOs(),
		TestType:    commConsts.TestUnit,
		TestTool:    req.TestTool,
		BuildTool:   req.BuildTool,
		TestCommand: req.Cmd,

		StartTime: startTime,
		EndTime:   endTime,
		Pass:      0, Fail: 0, Total: 0,

		SubmitResult:  req.SubmitResult,
		WorkspaceId:   req.WorkspaceId,
		WorkspaceType: req.WorkspaceType,
	}

	failedCaseLines, failedCaseLinesDesc := GenUnitReport(cases, &report, duration)

	// print case result one by one
	width := strconv.Itoa(len(strconv.Itoa(report.Total)))
	titleMaxWidth := getTitleMaxWidth(cases)
	for idx, cs := range cases {
		testSuite := stringUtils.AddPostfix(cs.TestSuite, classNameMaxWidth, commConsts.SpaceQuote)

		csTitle := cs.Title
		lent := runewidth.StringWidth(cs.Title)
		if titleMaxWidth > lent {
			postFix := strings.Repeat(commConsts.SpaceQuote, titleMaxWidth-lent)
			csTitle += postFix
		}

		status := GenStatusTxt(cs.Status)

		format := "(%" + width + "d/%d) [%s] [%s] [%s] [%.3fs]"
		msgCase := fmt.Sprintf(format, idx+1, report.Total, status, testSuite, csTitle, cs.Duration)

		websocketHelper.SendExecMsgIfNeed(msgCase, "", commConsts.Result, nil, wsMsg)

		logUtils.ExecConsolef(-1, msgCase)
		logUtils.ExecResult(msgCase)
	}

	logUtils.ExecConsolef(-1, "")

	// print failed cases with whole result status
	msg := ""
	status := commConsts.PASS
	msgCategory := commConsts.Result
	if report.Fail > 0 {
		status = commConsts.FAIL
		msgCategory = commConsts.Error
		msg = genUnitFailedMsg(failedCaseLines, failedCaseLinesDesc)

		logUtils.ExecConsolef(-1, msg)
		logUtils.ExecResult(msg)
	}

	websocketHelper.SendExecMsgIfNeed(msg, "", msgCategory,
		iris.Map{"key": key, "status": status}, wsMsg)

	// print summary result
	// 执行%d个用例，耗时%d秒%s。%s，%s，%s。报告%s。
	runResult, msgRunColor := GenRunResult(report)

	msgRun := dateUtils.DateTimeStr(time.Now()) + " " + runResult

	websocketHelper.SendExecMsgIfNeed(msgRunColor, "", commConsts.Result, nil, wsMsg)
	logUtils.ExecResult(msgRun)

	// print result path
	resultPath := filepath.Join(commConsts.ExecLogDir, commConsts.ResultText)
	msgReport := i118Utils.Sprintf("run_report") + " " + resultPath + "."
	if commConsts.ExecFrom == commConsts.FromCmd {
		msgReport = color.New(color.Bold, color.FgHiWhite).Sprint(i118Utils.Sprintf("run_report")) + " [" + resultPath + "]"
	}

	websocketHelper.SendExecMsgIfNeed(msgReport, "false", commConsts.Result, map[string]interface{}{
		"logDir": commConsts.ExecLogDir,
	}, wsMsg)

	logUtils.ExecConsole(-1, msgReport)
	logUtils.ExecConsole(-1, msgRun+"\n")
	logUtils.ExecResult(msgReport)
	report.Log = fileUtils.ReadFile(filepath.Join(commConsts.ExecLogDir, commConsts.LogText))

	//report.ProductId, _ = strconv.Atoi(vari.ProductId)
	json, _ := json.MarshalIndent(report, "", "\t")
	jsonPath := filepath.Join(commConsts.ExecLogDir, commConsts.ResultJson)
	fileUtils.WriteFile(jsonPath, string(json))

	return
}

func genUnitFailedMsg(failedCaseLines []string, failedCaseLinesDesc []string) string {
	divider := shellUtils.GenFullScreenDivider()

	msg := divider
	msg += "\n" + color.New(color.Bold, color.FgHiWhite).Sprint(i118Utils.Sprintf("failed_scripts")) + "\n"
	msg += strings.Join(failedCaseLines, "\n")
	msg += strings.Join(failedCaseLinesDesc, "\n")
	msg += "\n\n" + divider

	return msg
}

func getTitleMaxWidth(cases []commDomain.UnitResult) int {
	maxWidth := 0
	for _, cs := range cases {
		titleWidth := runewidth.StringWidth(cs.Title)
		if maxWidth < titleWidth {
			maxWidth = titleWidth
		}
	}

	return maxWidth
}

func RetrieveUnitResult(testset serverDomain.TestSet, startTime int64) (
	suites []commDomain.UnitTestSuite) {

	resultFiles := make([]string, 0)
	if testset.ResultDir != "" {
		resultFiles, _ = GetSuiteFiles(testset.ResultDir, startTime, testset.TestTool)
	}

	if isAllureReport(testset.TestTool) {
		if testset.ResultDir != "" {
			suites = GetAllureSuites(testset.ResultDir, startTime)
		} else {
			logUtils.Info(color.RedString(
				i118Utils.Sprintf("must_provide_allure_report_dir")))
		}

	} else {
		failedCaseIdToThresholdMap := map[string]string{}

		if testset.TestTool == commConsts.K6 && len(resultFiles) > 1 {
			content := fileUtils.ReadFile(resultFiles[len(resultFiles)-1])
			failedCaseIdToThresholdMap = GetK6FailCaseInSummary(content)
		}

		for _, file := range resultFiles {
			testSuite, err := GetTestSuite(file, testset.TestTool, failedCaseIdToThresholdMap)

			if err == nil {
				suites = append(suites, testSuite)
			}
		}
	}

	return
}

func GetAllureSuites(resultDir string, startTime int64) (suites []commDomain.UnitTestSuite) {
	files, err := ioutil.ReadDir(resultDir)
	if err != nil {
		return
	}

	cases := make([]commDomain.AllureCase, 0)
	for _, fi := range files {
		name := fi.Name()

		if strings.Index(name, "-result.json") < 0 { // || fi.ModTime().Unix() < startTime {
			continue
		}

		pth := filepath.Join(resultDir, name)
		content := fileUtils.ReadFileBuf(pth)

		cs := commDomain.AllureCase{}
		err = json.Unmarshal(content, &cs)
		if err == nil {
			cases = append(cases, cs)
		}
	}

	suites = ConvertAllureResult(cases)

	return
}

func GetSuiteFiles(resultDir string, startTime int64, testTool commConsts.TestTool) (resultFiles []string, err error) {
	if fileUtils.IsDir(resultDir) {
		dir, err := ioutil.ReadDir(resultDir)
		if err == nil {
			for _, fi := range dir {
				name := fi.Name()
				ext := path.Ext(name)

				//if fi.ModTime().Unix() < startTime {
				//	continue
				//}

				if ((isAllureReport(testTool) || testTool == commConsts.K6) && ext == ".json") || ext == ".xml" {
					pth := filepath.Join(resultDir, name)
					resultFiles = append(resultFiles, pth)
				}
			}
		}
	} else {
		resultFiles = append(resultFiles, resultDir)
	}

	return
}

func GetTestSuite(logFile string, testTool commConsts.TestTool, failedCaseIdToThresholdMap map[string]string) (
	testSuite commDomain.UnitTestSuite, err error) {

	content := fileUtils.ReadFile(logFile)

	if testTool == commConsts.JUnit || testTool == commConsts.TestNG {
		testSuite = commDomain.UnitTestSuite{}
		err = xml.Unmarshal([]byte(content), &testSuite)

	} else if testTool == commConsts.PHPUnit {
		phpTestSuite := commDomain.PhpUnitSuites{}
		err = xml.Unmarshal([]byte(content), &phpTestSuite)
		if err == nil {
			testSuite = ConvertPhpUnitResult(phpTestSuite)
		}
	} else if testTool == commConsts.PyTest {
		pyTestSuite := commDomain.PyTestSuites{}
		err = xml.Unmarshal([]byte(content), &pyTestSuite)
		if err == nil {
			testSuite = ConvertPyTestResult(pyTestSuite)
		}
	} else if testTool == commConsts.Jest {
		jestSuite := commDomain.JestSuites{}
		err = xml.Unmarshal([]byte(content), &jestSuite)
		if err == nil {
			testSuite = ConvertJestResult(jestSuite)
		}
	} else if testTool == commConsts.GTest {
		gTestSuite := commDomain.GTestSuites{}
		err = xml.Unmarshal([]byte(content), &gTestSuite)
		if err == nil {
			testSuite = ConvertGTestResult(gTestSuite)
		}
	} else if testTool == commConsts.QTest {
		qTestSuite := commDomain.QTestSuites{}
		err = xml.Unmarshal([]byte(content), &qTestSuite)
		if err == nil {
			testSuite = ConvertQTestResult(qTestSuite)
		}
	} else if testTool == commConsts.CppUnit {
		content = strings.Replace(content, "ISO-8859-1", "UTF-8", -1)

		cppUnitSuites := commDomain.CppUnitSuites{}
		err = xml.Unmarshal([]byte(content), &cppUnitSuites)
		if err == nil {
			testSuite = ConvertCppUnitResult(cppUnitSuites)
		}
	} else if testTool == commConsts.RobotFramework {
		robotResult := commDomain.RobotResult{}
		err = xml.Unmarshal([]byte(content), &robotResult)
		if err == nil {
			testSuite = ConvertRobotResult(robotResult)
		}
	} else if testTool == commConsts.Cypress || testTool == commConsts.Playwright {
		result := commDomain.CypressTestsuites{}
		err = xml.Unmarshal([]byte(content), &result)
		if err == nil {
			testSuite = ConvertCyResult(result)
		}
	} else if testTool == commConsts.Puppeteer {
		cyResult := commDomain.CypressTestsuites{}

		cySuite := commDomain.CypressTestsuite{}
		err = xml.Unmarshal([]byte(content), &cySuite)

		cyResult.Testsuites = append(cyResult.Testsuites, cySuite)
		if err == nil {
			testSuite = ConvertCyResult(cyResult)
		}
	} else if testTool == commConsts.K6 {
		results := []interface{}{}
		lines := strings.Split(content, "\n")

		if len(lines) > 1 {
			for _, line := range strings.Split(content, "\n") {
				k6Point := commDomain.K6Point{}
				errInner := json.Unmarshal([]byte(line), &k6Point)
				if errInner == nil && k6Point.Type == commConsts.Point {
					results = append(results, k6Point)
					continue
				}
			}
			testSuite = ConvertK6Result(results, failedCaseIdToThresholdMap)
		}
	}

	return
}

func ParserUnitTestResult(testSuites []commDomain.UnitTestSuite) (
	cases []commDomain.UnitResult, classNameMaxWidth int, dur float32) {

	idx := 1
	for _, suite := range testSuites {
		for _, cs := range suite.Cases {
			getCaseIdFromName(&cs, idx)

			if cs.Failure != nil {
				cs.Status = "fail"

				cs.Failure.Desc = strings.Replace(cs.Failure.Desc, "<![CDATA[", "", -1)
				cs.Failure.Desc = strings.Replace(cs.Failure.Desc, "]]>", "", -1)
			} else if cs.ErrorContent != "" {
				cs.Status = "fail"

				if cs.Failure == nil {
					cs.Failure = &commDomain.Failure{}
				}
				cs.ErrorContent = strings.Replace(cs.ErrorContent, "<![CDATA[", "", -1)
				cs.ErrorContent = strings.Replace(cs.ErrorContent, "]]>", "", -1)
				cs.Failure.Desc = cs.ErrorType + ": " + cs.ErrorContent
			} else {
				cs.Status = "pass"
			}

			lent2 := runewidth.StringWidth(cs.TestSuite)
			if lent2 > classNameMaxWidth {
				classNameMaxWidth = lent2
			}

			cases = append(cases, cs)
			idx++
		}
	}

	return
}

func ConvertAllureResult(cases []commDomain.AllureCase) (testSuites []commDomain.UnitTestSuite) {
	suites := make([]*commDomain.UnitTestSuite, 0)
	suiteMap := map[string]*commDomain.UnitTestSuite{}

	for _, cs := range cases {
		suiteName := GetAllureCaseSuiteName(cs)
		//logUtils.Info(suiteName)

		_, ok := suiteMap[suiteName]
		if !ok {
			suite := commDomain.UnitTestSuite{
				Name: suiteName,
				Time: 0,
			}
			suites = append(suites, &suite)
			suiteMap[suiteName] = &suite
		}

		//suiteMap[suiteName].Name = "111"

		caseId := GetAllureCaseId(cs)

		// passed, failed
		var status commConsts.ResultStatus
		if cs.Status == "passed" {
			status = commConsts.PASS
		} else if cs.Status == "failed" {
			status = commConsts.FAIL
		}
		caseResult := commDomain.UnitResult{
			Id:        caseId,
			Cid:       caseId,
			Title:     cs.Name,
			TestSuite: suiteName,
			Duration:  float32(cs.Stop-cs.Start) / 1000,
			Status:    status,
		}

		if cs.Status == "failed" {
			caseResult.Failure = &commDomain.Failure{
				Type: "AssertionError",
				Desc: cs.StatusDetails.Message + ": " + cs.StatusDetails.Trace,
			}
		}

		suiteMap[suiteName].Cases = append(suiteMap[suiteName].Cases, caseResult)
	}

	for _, suite := range suites {
		sort.Sort(commDomain.CaseSlice(suite.Cases))

		dur := int64(0)
		for _, cs := range suite.Cases {
			dur += cs.EndTime - cs.StartTime
		}
		suite.Time = float32(dur)

		testSuites = append(testSuites, *suite)
	}

	return
}

func GetAllureCaseSuiteName(cs commDomain.AllureCase) (name string) {
	suiteArr := make([]string, 0)

	for _, label := range cs.Labels {
		if label.Name == "parentSuite" {
			if label.Value != "" {
				suiteArr = append(suiteArr, label.Value)
			}
		} else if label.Name == "suite" || label.Name == "subSuite" {
			if label.Value != "" && (len(suiteArr) == 0 || label.Value != suiteArr[len(suiteArr)-1]) {
				suiteArr = append(suiteArr, label.Value)
			}
		}
	}

	name = strings.Join(suiteArr, "-")

	return
}

func GetAllureCaseId(cs commDomain.AllureCase) (id int) {
	// 1. from testCaseId
	id, err := strconv.Atoi(cs.TestCaseId)
	if err == nil && id > 0 {
		return
	}

	// 2. from as_id label
	for _, label := range cs.Labels {
		if label.Name == "as_id" {
			if label.Value != "" {
				cs.TestCaseId = label.Value // 2
			}

			break
		}
	}
	id = stringUtils.ParseInt(cs.TestCaseId)

	// 2. from the ids param in name like [cs-1]
	regx := regexp.MustCompile(`\[(.+)\]`)
	arr := regx.FindAllStringSubmatch(cs.Name, -1)
	if len(arr) > 0 {
		item := arr[len(arr)-1]
		idFromName := stringUtils.ParseInt(item[1])
		if idFromName > 0 {
			id = idFromName
		}
	}

	return
}

func ConvertJestResult(jestSuite commDomain.JestSuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}
	testSuite.Time = jestSuite.Time

	for _, suite := range jestSuite.TestSuites {
		for _, cs := range testSuite.Cases {
			caseResult := commDomain.UnitResult{}
			caseResult.Title = cs.Title
			caseResult.Duration = cs.Duration

			if suite.Title != "" && suite.Title != "undefined" {
				caseResult.TestSuite = suite.Title
			} else {
				caseResult.TestSuite = jestSuite.Title
			}

			caseResult.Failure = cs.Failure

			testSuite.Cases = append(testSuite.Cases, caseResult)
		}
	}

	return testSuite
}

func ConvertPhpUnitResult(phpUnitSuite commDomain.PhpUnitSuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}

	var total float32 = 0
	for _, cs := range phpUnitSuite.Cases {
		caseResult := commDomain.UnitResult{}
		caseResult.Title = cs.Title
		caseResult.Duration = cs.Time

		total += cs.Time

		if cs.Groups != "" && cs.Groups != "default" {
			caseResult.TestSuite = cs.Groups
		} else {
			caseResult.TestSuite = cs.TestSuite
		}

		if cs.Status != 0 {
			fail := commDomain.Failure{}
			fail.Desc = cs.Fail
			caseResult.Failure = &fail
		}

		testSuite.Cases = append(testSuite.Cases, caseResult)
	}
	testSuite.Duration = int64(total)
	testSuite.Time = total

	return testSuite
}

func ConvertPyTestResult(pytestSuites commDomain.PyTestSuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}

	var total float32 = 0
	for _, suite := range pytestSuites.TestSuites {
		total += suite.Time

		for _, cs := range suite.Cases {
			caseResult := commDomain.UnitResult{}
			caseResult.Title = cs.Title
			caseResult.Duration = cs.Duration

			if suite.Title != "" && suite.Title != "pytest" {
				caseResult.TestSuite = suite.Title
			} else {
				caseResult.TestSuite = cs.TestSuite
			}

			if cs.Failure != nil {
				fail := commDomain.Failure{}
				fail.Type = cs.Failure.Type
				fail.Desc = cs.Failure.Desc
				caseResult.Failure = &fail
			} else if cs.Error != nil {
				fail := commDomain.Failure{}
				fail.Type = cs.Error.Message
				fail.Desc = cs.Error.Text
				caseResult.Failure = &fail
			}

			testSuite.Cases = append(testSuite.Cases, caseResult)

		}
	}

	testSuite.Duration = int64(total)
	testSuite.Time = total

	return testSuite
}

func ConvertGTestResult(gTestSuite commDomain.GTestSuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}
	testSuite.Time = gTestSuite.Time

	for _, suite := range gTestSuite.TestSuites {
		for _, cs := range suite.Cases {
			caseResult := commDomain.UnitResult{}
			caseResult.Title = cs.Title
			caseResult.Duration = cs.Duration
			caseResult.Status = commConsts.ResultStatus(cs.Status)

			if suite.Title != "" && suite.Title != "pytest" {
				caseResult.TestSuite = suite.Title
			}

			if cs.Failure != nil {
				fail := commDomain.Failure{}
				fail.Type = cs.Failure.Type
				fail.Desc = cs.Failure.Desc
				caseResult.Failure = &fail
			}

			testSuite.Cases = append(testSuite.Cases, caseResult)

		}
	}

	return testSuite
}

func ConvertCppUnitResult(cppunitSuite commDomain.CppUnitSuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}

	for _, cs := range cppunitSuite.FailedTests.Cases {
		caseResult := commDomain.UnitResult{}
		caseResult.Id = cs.Id
		caseResult.Title = cs.Title

		fail := commDomain.Failure{}
		fail.Type = cs.FailureType
		fail.Desc = cs.Message
		caseResult.Failure = &fail

		testSuite.Cases = append(testSuite.Cases, caseResult)
	}

	for _, cs := range cppunitSuite.SuccessfulTests.Cases {
		caseResult := commDomain.UnitResult{}
		caseResult.Id = cs.Id
		caseResult.Title = cs.Title

		testSuite.Cases = append(testSuite.Cases, caseResult)
	}

	return testSuite
}

func ConvertQTestResult(qTestSuite commDomain.QTestSuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}

	for _, cs := range qTestSuite.Cases {
		caseResult := commDomain.UnitResult{}
		caseResult.TestSuite = qTestSuite.Name
		caseResult.Title = cs.Title
		caseResult.Status = commConsts.ResultStatus(cs.Result)

		if cs.Failure != nil {
			fail := commDomain.Failure{}
			fail.Type = cs.Failure.Type
			fail.Desc = cs.Failure.Desc
			caseResult.Failure = &fail
		}

		testSuite.Cases = append(testSuite.Cases, caseResult)
	}

	return testSuite
}

func ConvertRobotResult(result commDomain.RobotResult) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}

	suiteMap := map[string]string{}
	for _, state := range result.Statistics.Suite.States {
		suiteMap[state.ID] = state.Text
	}

	tests := make([]commDomain.RobotTest, 0)
	for _, suite := range result.Suites {
		RetrieveRobotTests(suite, &tests)
	}

	for _, cs := range tests {
		caseResult := commDomain.UnitResult{}
		caseResult.Title = cs.Name
		caseResult.Status = commConsts.ResultStatus(strings.ToLower(cs.Status.Status))

		suiteId := cs.ID[0:strings.LastIndex(cs.ID, "-")]
		caseResult.TestSuite = suiteMap[suiteId]

		templ := "20060102 15:04:05.000"
		startTime, _ := time.ParseInLocation(templ, cs.Status.StartTime, time.Local)
		endTime, _ := time.ParseInLocation(templ, cs.Status.EndTime, time.Local)

		caseResult.StartTime = startTime.Unix()
		caseResult.EndTime = endTime.Unix()
		caseResult.Duration = float32(caseResult.EndTime - caseResult.StartTime)

		if caseResult.Status != "pass" {
			fail := commDomain.Failure{}
			fail.Type = ""
			fail.Desc = cs.Status.Text
			caseResult.Failure = &fail
		}

		testSuite.Cases = append(testSuite.Cases, caseResult)
	}

	return testSuite
}

func RetrieveRobotTests(suite commDomain.RobotSuite, tests *[]commDomain.RobotTest) {
	for _, suite := range suite.Suites {
		RetrieveRobotTests(suite, tests)
	}

	for _, test := range suite.Tests {
		*tests = append(*tests, test)
	}
}

func ConvertCyResult(result commDomain.CypressTestsuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}

	for _, suite := range result.Testsuites {
		if suite.Name == "Root Suite" {
			continue
		}

		templ := "20060102 15:04:05.000"
		duration := suite.Time
		startTime, err := time.ParseInLocation(templ, suite.Timestamp, time.Local)
		//endTime := time.Unix(startTime.Unix() + int64(duration), 0)

		if err != nil {
			startTime, err = time.ParseInLocation(time.RFC1123, suite.Timestamp, time.Local)
		}

		testSuite.Duration = int64(duration)
		testSuite.Time = float32(startTime.Unix())

		for _, cs := range suite.Testcases {
			caseResult := commDomain.UnitResult{}
			caseResult.TestSuite = suite.Name
			caseResult.Title = cs.Name
			caseResult.Duration = float32(cs.Time)

			if len(cs.Failures) > 0 {
				caseResult.Status = "fail"

				fail := commDomain.Failure{}
				fail.Type = cs.Failures[0].Type
				fail.Desc = cs.Failures[0].Message
				caseResult.Failure = &fail
			} else {
				caseResult.Status = "pass"
			}

			testSuite.Cases = append(testSuite.Cases, caseResult)
		}
	}

	return testSuite
}

func isAllureReport(testTool commConsts.TestTool) (ret bool) {
	return testTool == commConsts.Allure || testTool == commConsts.GoTest
}

func getResultDirForDifferentTool(testset *serverDomain.TestSet) {
	if testset.TestTool == commConsts.JUnit && testset.BuildTool == commConsts.Maven {
		testset.ResultDir = filepath.Join("target", "surefire-reports")
		testset.ZipDir = testset.ResultDir

	} else if testset.TestTool == commConsts.TestNG && testset.BuildTool == commConsts.Maven {
		testset.ResultDir = filepath.Join("target", "surefire-reports", "junitreports")
		testset.ZipDir = filepath.Dir(testset.ResultDir)

	} else if testset.TestTool == commConsts.RobotFramework || testset.TestTool == commConsts.Cypress ||
		testset.TestTool == commConsts.Playwright || testset.TestTool == commConsts.Puppeteer ||
		testset.TestTool == commConsts.K6 {
		testset.ResultDir = "results"
		testset.ZipDir = testset.ResultDir
	} else if testset.TestTool == commConsts.Zap {
		testset.ResultDir = getZapReport()
		testset.ZipDir = testset.ResultDir

	} else if isAllureReport(testset.TestTool) {
		testset.ResultDir = commConsts.AllureReportDir
		testset.ZipDir = testset.ResultDir

	} else {
		testset.ResultDir = "testresults.xml"
		testset.ZipDir = testset.ResultDir
	}

	if testset.ResultDir != "" {
		if !fileUtils.IsAbsolutePath(testset.ResultDir) {
			testset.ResultDir = filepath.Join(testset.WorkspacePath, testset.ResultDir)
		}

		if !fileUtils.IsAbsolutePath(testset.ZipDir) {
			testset.ZipDir = filepath.Join(testset.WorkspacePath, testset.ZipDir)
		}
	}

	return
}

func getZapReport() (ret string) {
	for index, item := range os.Args {
		if item == "-quickout" && index < len(os.Args)-1 {
			ret = os.Args[index+1]
			return
		}
	}

	return
}

func getCaseIdFromName(cs *commDomain.UnitResult, defaultVal int) {
	if cs.Cid > 0 {
		return
	}

	cs.Title = strings.TrimSpace(cs.Title)

	regx := regexp.MustCompile(`^(\d+)\. (.+)`)
	arr := regx.FindAllStringSubmatch(cs.Title, -1)
	if len(arr) > 0 {
		cs.Cid = stringUtils.ParseInt(arr[0][1])
		cs.Title = arr[0][2]
	}

	cs.Id = cs.Cid
	if cs.Id <= 0 {
		cs.Id = defaultVal
	}

	return
}
