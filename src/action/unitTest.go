package action

import (
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	shellUtils "github.com/easysoft/zentaoatf/src/utils/shell"
	"time"
)

func RunUnitTest(cmdStr string) string {
	startTime := time.Now().Unix()
	shellUtils.ExeAppWithOutput(cmdStr)
	endTime := time.Now().Unix()

	testSuites, resultDir := testingService.RetrieveUnitResult(startTime)
	cases, classNameMaxWidth, duration := testingService.ParserUnitTestResult(testSuites)

	if duration == 0 {
		duration = float32(endTime - startTime)
	}

	report := testingService.GenUnitTestReport(cases, classNameMaxWidth, duration)

	zentaoService.CommitTestResult(report, 0)

	return resultDir
}
