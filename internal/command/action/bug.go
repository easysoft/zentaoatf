package action

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	analysisUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/analysis"
	"github.com/aaronchen2k/deeptest/internal/command"
	"github.com/aaronchen2k/deeptest/internal/command/ui/page"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/stdin"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
)

func CommitBug(files []string, actionModule *command.IndexModule) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		stdinUtils.InputForDir(&resultDir, "", "result")
	}
	//resultDir = fileUtils.AddPathSepIfNeeded(resultDir)

	report, err := analysisUtils.ReadReport(commConsts.WorkDir, resultDir)
	if err != nil {
		return
	}

	ids := make([]string, 0)
	lines := make([]string, 0)
	for _, cs := range report.FuncResult {
		if cs.Status != commConsts.PASS {
			lines = append(lines, fmt.Sprintf("%d. %s %s", cs.Id, cs.Title, coloredStatus(cs.Status)))
			ids = append(ids, strconv.Itoa(cs.Id))
		}
	}

	for {
		logUtils.ExecConsole(color.FgCyan, "\n"+i118Utils.Sprintf("enter_case_id_for_report_bug"))
		logUtils.ExecConsole(color.FgCyan, strings.Join(lines, "\n"))
		var caseId string
		fmt.Scanln(&caseId)
		if caseId == "exit" {
			color.Unset()
			os.Exit(0)
		} else {
			if stringUtils.FindInArr(caseId, ids) {
				page.CuiReportBug(resultDir, caseId, actionModule)
			} else {
				logUtils.ExecConsole(color.FgRed, i118Utils.Sprintf("invalid_input"))
			}
		}
	}
}

func coloredStatus(status commConsts.ResultStatus) string {
	temp := strings.ToLower(status.String())

	switch temp {
	case "pass":
		return color.GreenString(i118Utils.Sprintf(temp))
	case "fail":
		return color.RedString(i118Utils.Sprintf(temp))
	case "skip":
		return color.YellowString(i118Utils.Sprintf(temp))
	}

	return status.String()
}
