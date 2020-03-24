package action

import (
	"github.com/easysoft/zentaoatf/src/model"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	shellUtils "github.com/easysoft/zentaoatf/src/utils/shell"
	"github.com/mattn/go-runewidth"
	"strings"
)

func UnitTest(cmdStr string) {
	logUtils.InitLogger()

	shellUtils.ExeShellWithOutput(cmdStr)
	testResult := testingService.RetriveResult()

	classNameMaxWidth := 0

	count := 1
	var report = model.TestReport{Env: commonUtils.GetOs(), Pass: 0, Fail: 0, Total: 0}
	for idx, cs := range testResult.Testcase {
		testResult.Testcase[idx].Id = count
		if cs.Failure != nil {
			testResult.Testcase[idx].Status = "fail"

			cs.Failure.Desc = strings.Replace(cs.Failure.Desc, "<![CDATA[", "", -1)
			cs.Failure.Desc = strings.Replace(cs.Failure.Desc, "]]>", "", -1)
			logUtils.Screen(cs.Failure.Desc)
		} else {
			testResult.Testcase[idx].Status = "pass"
		}

		lent2 := runewidth.StringWidth(cs.Classname)
		if lent2 > classNameMaxWidth {
			classNameMaxWidth = lent2
		}

		count++
	}

	testingService.UnitTestReport(report, testResult, classNameMaxWidth)

}
