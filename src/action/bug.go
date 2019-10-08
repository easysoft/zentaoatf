package action

import (
	"fmt"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	"github.com/easysoft/zentaoatf/src/ui/page"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	stdinUtils "github.com/easysoft/zentaoatf/src/utils/stdin"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
)

func CommitBug(files []string) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		stdinUtils.InputForDir(&resultDir, "", "result")
	}
	resultDir = fileUtils.UpdateDir(resultDir)

	report := testingService.GetTestTestReportForSubmit(resultDir)

	ids := make([]string, 0)
	lines := make([]string, 0)
	for _, cs := range report.Cases {
		if cs.Status != constant.PASS.String() {
			lines = append(lines, fmt.Sprintf("%d. %s %s", cs.Id, cs.Title, logUtils.ColoredStatus(cs.Status)))
			ids = append(ids, strconv.Itoa(cs.Id))
		}
	}

	for {
		logUtils.PrintToStdOut("\n"+i118Utils.I118Prt.Sprintf("enter_case_id_for_report_bug"), color.FgCyan)
		logUtils.PrintToStdOut(strings.Join(lines, "\n"), -1)

		var caseId string
		fmt.Scanln(&caseId)

		if caseId == "exit" {
			os.Exit(1)
		} else {
			if stringUtils.FindInArr(caseId, ids) {
				page.CuiReportBug(resultDir, caseId)
			} else {
				logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("invalid_input"), color.FgRed)
			}
		}
	}
}
