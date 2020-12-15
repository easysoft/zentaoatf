package action

import (
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	shellUtils "github.com/easysoft/zentaoatf/src/utils/shell"
)

func RunUnitTest(cmdStr string) {
	shellUtils.ExeShellWithOutput(cmdStr)

	testSuites := testingService.RetrieveUnitResult()
	cases, classNameMaxWidth, time := testingService.ParserUnitTestResult(testSuites)

	report := testingService.GenUnitTestReport(cases, classNameMaxWidth, time)

	zentaoService.CommitTestResult(report, 0)
}
