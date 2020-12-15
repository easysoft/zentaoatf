package action

import (
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	shellUtils "github.com/easysoft/zentaoatf/src/utils/shell"
	time2 "time"
)

func RunUnitTest(cmdStr string) {
	startTime := time2.Now().Unix()
	shellUtils.ExeShellWithOutput(cmdStr)
	endTime := time2.Now().Unix()

	testSuites := testingService.RetrieveUnitResult()
	cases, classNameMaxWidth, time := testingService.ParserUnitTestResult(testSuites)

	if time == 0 {
		time = float32(endTime - startTime)
	}

	report := testingService.GenUnitTestReport(cases, classNameMaxWidth, time)

	zentaoService.CommitTestResult(report, 0)
}
