package action

import (
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	shellUtils "github.com/easysoft/zentaoatf/src/utils/shell"
	"time"
)

func UnitTest(cmdStr string) {
	startTime := time.Now().Unix()
	shellUtils.ExeShellWithOutput(cmdStr)
	endTime := time.Now().Unix()

	testSuites := testingService.RetriveResult()
	cases, classNameMaxWidth := testingService.ParserUnitTestResult(testSuites)

	report := testingService.GenUnitTestReport(cases, classNameMaxWidth)

	report.StartTime = startTime
	report.EndTime = endTime
	report.Duration = endTime - startTime
	zentaoService.CommitUnitTestResult(report)
}
