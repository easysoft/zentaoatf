package action

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	analysisUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/analysis"
	configHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	"github.com/aaronchen2k/deeptest/internal/command"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/stdin"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
)

var (
	bug       commDomain.ZtfBug
	bugFields commDomain.ZentaoBugFields
)

func CommitBug(files []string, actionModule *command.IndexModule) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		stdinUtils.InputForDir(&resultDir, "", "result")
	}
	//resultDir = fileUtils.AddPathSepIfNeeded(resultDir)

	report, err := analysisUtils.ReadReportByWorkspaceSeq(commConsts.WorkDir, resultDir)
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
				reportBug(resultDir, caseId)
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

func reportBug(resultDir string, caseId string) error {
	bug = zentaoUtils.PrepareBug(commConsts.WorkDir, resultDir, caseId)

	config := configHelper.LoadByWorkspacePath(commConsts.WorkDir)

	bugFields, _ = zentaoUtils.GetBugFiledOptions(config, bug.Product)

	//bug.Module = 0
	bug.Severity = 3
	bug.Pri = 3
	bug.Type = getFirstNoEmptyVal(bugFields.Types)
	bug.Version = getNameById(bug.Version, bugFields.Build)

	err := zentaoUtils.CommitBug(bug, config)
	return err
}

func getFirstNoEmptyVal(options []commDomain.BugOption) string {
	for _, opt := range options {
		if opt.Name != "" {
			return opt.Code
		}
	}

	return ""
}

func getNameById(id string, options []commDomain.BugOption) string {
	for _, opt := range options {
		if opt.Code == id {
			return opt.Name
		}
	}

	return ""
}
