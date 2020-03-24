package action

import (
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	shellUtils "github.com/easysoft/zentaoatf/src/utils/shell"
)

func UnitTest(cmdStr string) {
	shellUtils.ExeShellWithOutput(cmdStr)

	testSuites := testingService.RetriveResult()
	cases, classNameMaxWidth := testingService.ParserUnitTestResult(testSuites)

	report := testingService.GenUnitTestReport(cases, classNameMaxWidth)
	testingService.SubmitUnitTestReport(report)
}
