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
	resultDir = fileUtils.AddPathSepIfNeeded(resultDir)

	report := testingService.GetZTFTestReportForSubmit(resultDir)

	ids := make([]string, 0)
	lines := make([]string, 0)
	for _, cs := range report.FuncResult {
		if cs.Status != constant.PASS.String() {
			lines = append(lines, fmt.Sprintf("%d. %s %s", cs.Id, cs.Title, logUtils.ColoredStatus(cs.Status)))
			ids = append(ids, strconv.Itoa(cs.Id))
		}
	}

	if len(ids) == 0 {
		logUtils.PrintToWithColor("\n"+i118Utils.Sprintf("no_fail_cases"), color.FgCyan)
		return
	}

	for {
		logUtils.PrintToWithColor("\n"+i118Utils.Sprintf("enter_case_id_for_report_bug"), color.FgCyan)
		logUtils.PrintToWithColor(strings.Join(lines, "\n"), -1)

		var caseId string
		fmt.Scanln(&caseId)

		if caseId == "exit" {
			color.Unset()
			os.Exit(0)
		} else {
			if stringUtils.FindInArr(caseId, ids) {
				page.CuiReportBug(resultDir, caseId)
			} else {
				logUtils.PrintToWithColor(i118Utils.Sprintf("invalid_input"), color.FgRed)
			}
		}
	}
}
