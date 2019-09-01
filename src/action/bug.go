package action

import (
	"fmt"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	"github.com/easysoft/zentaoatf/src/ui/page"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/fatih/color"
	"os"
	"regexp"
)

func CommitBug(files []string) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		configUtils.ConfigForDir(&resultDir, "result")
	}
	resultDir = fileUtils.UpdateDir(resultDir)

	report := testingService.GetTestTestReportForSubmit(resultDir)

	for {
		var caseId string

		logUtils.PrintToStdOut(i118Utils.I118Prt.Sprint("enter_case_id_for_report_bug"), color.FgCyan)
		for _, cs := range report.Cases {
			if cs.Status != constant.PASS.String() {
				logUtils.PrintToStdOut(fmt.Sprintf("\n%d. %s %s", cs.Id, cs.Title, logUtils.ColoredStatus(cs.Status)),
					color.FgCyan)
			}
		}
		fmt.Scanln(&caseId)

		if caseId == "exit" {
			os.Exit(1)
		} else {
			pass, _ := regexp.MatchString("^\\d+$", caseId)
			if pass {
				page.CuiReportBug(resultDir, caseId)
			} else {
				logUtils.PrintToStdOut(i118Utils.I118Prt.Sprint("invalid_input"), color.FgRed)
			}
		}
	}
}
