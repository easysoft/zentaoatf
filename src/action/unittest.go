package action

import (
	"github.com/easysoft/zentaoatf/src/model"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	shellUtils "github.com/easysoft/zentaoatf/src/utils/shell"
)

func UnitTest(cmdStr string) string {
	logUtils.InitLogger()

	shellUtils.ExeShellWithOutput(cmdStr)
	testResult := testingService.RetriveResult()

	pathMaxWidth := 0
	//numbMaxWidth := 0
	//for _, suite := range testResult.Suites {
	//	for _, cs := range suite.Cases {
	//
	//	}
	//}

	var report = model.TestReport{Env: commonUtils.GetOs(), Pass: 0, Fail: 0, Total: 0}

	testingService.UnitTestReport(report, testResult, pathMaxWidth)

	return ""
}
