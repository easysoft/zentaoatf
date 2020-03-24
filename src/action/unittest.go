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
	testSuites := testingService.RetriveResult()

	cases := make([]model.UnitTestCase, 0)
	classNameMaxWidth := 0
	idx := 1
	for _, suite := range testSuites {
		for _, cs := range suite.Testcase {
			cs.Id = idx
			if cs.Failure != nil {
				cs.Status = "fail"

				cs.Failure.Desc = strings.Replace(cs.Failure.Desc, "<![CDATA[", "", -1)
				cs.Failure.Desc = strings.Replace(cs.Failure.Desc, "]]>", "", -1)
				logUtils.Screen(cs.Failure.Desc)
			} else {
				cs.Status = "pass"
			}

			lent2 := runewidth.StringWidth(cs.Classname)
			if lent2 > classNameMaxWidth {
				classNameMaxWidth = lent2
			}

			cases = append(cases, cs)
			idx++
		}
	}

	var report = model.TestReport{Env: commonUtils.GetOs(), Pass: 0, Fail: 0, Total: 0}

	testingService.UnitTestReport(report, cases, classNameMaxWidth)

}
