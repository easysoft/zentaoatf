package action

import (
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	shellUtils "github.com/easysoft/zentaoatf/src/utils/shell"
	"time"
)

func RunUnitTest(cmdStr string) {
	startTime := time.Now().Unix()
	shellUtils.ExeShellWithOutput(cmdStr)
	endTime := time.Now().Unix()

	testSuites := testingService.RetrieveUnitResult()
	cases, classNameMaxWidth := testingService.ParserUnitTestResult(testSuites)

	report := testingService.GenUnitTestReport(cases, classNameMaxWidth, startTime, endTime)

	zentaoService.CommitTestResult(report)
}
