package action

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	analysisHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/analysis"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stdinUtils "github.com/easysoft/zentaoatf/pkg/lib/stdin"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"github.com/fatih/color"
)

var (
	bug       commDomain.ZtfBug
	bugFields commDomain.ZentaoBugFields
)

func CommitBug(files []string, productId int, noNeedConfirm bool) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		stdinUtils.InputForDir(&resultDir, "", "result")
	}

	if productId == 0 {
		productIdStr := stdinUtils.GetInput("\\d+", "",
			i118Utils.Sprintf("pls_enter")+" "+i118Utils.Sprintf("product_id"))
		productId, _ = strconv.Atoi(productIdStr)
	}

	report, pth, err := analysisHelper.ReadReportByWorkspaceSeq(commConsts.WorkDir, resultDir)
	if err != nil {
		logUtils.ExecConsole(color.FgRed, i118Utils.Sprintf("read_report_fail", pth))
		return
	}

	ids := make([]string, 0)
	lines := make([]string, 0)
	for _, cs := range report.FuncResult {
		if cs.Status != commConsts.PASS {
			lines = append(lines, fmt.Sprintf("%d. %s %s", cs.Id, cs.Title, coloredStatus(cs.Status.String())))
			ids = append(ids, strconv.Itoa(cs.Id))
		}
	}

	for _, cs := range report.UnitResult {
		if cs.Status != "pass" {
			lines = append(lines, fmt.Sprintf("%d. %s %s", cs.Id, cs.Title, coloredStatus(cs.Status.String())))
			ids = append(ids, strconv.Itoa(cs.Id))
		}
	}

	if len(lines) == 0 {
		logUtils.ExecConsole(color.FgCyan, i118Utils.Sprintf("no_failed_case_to_report_bug"))
		return
	}

	if noNeedConfirm {
		for _, caseId := range ids {
			reportBug(resultDir, caseId, productId)
		}

		return
	}

	// wait to input
	for {
		logUtils.ExecConsole(color.FgCyan, "\n"+i118Utils.Sprintf("enter_case_id_for_report_bug"))
		logUtils.ExecConsole(-1, strings.Join(lines, "\n"))
		var caseId string
		fmt.Scanln(&caseId)
		if caseId == "exit" {
			color.Unset()
			os.Exit(0)
		} else {
			if stringUtils.FindInArr(caseId, ids) {
				reportBug(resultDir, caseId, productId)
			} else {
				logUtils.ExecConsole(color.FgRed, i118Utils.Sprintf("invalid_input"))
			}
		}
	}
}

func coloredStatus(status string) string {
	temp := strings.ToLower(status)

	switch temp {
	case "pass":
		return color.GreenString(i118Utils.Sprintf(temp))
	case "fail":
		return color.RedString(i118Utils.Sprintf(temp))
	case "skip":
		return color.YellowString(i118Utils.Sprintf(temp))
	}

	return status
}

func reportBug(resultDir string, caseId string, productId int) (err error) {
	config := configHelper.LoadByWorkspacePath(commConsts.ZtfDir)

	bugFields, err = zentaoHelper.GetBugFiledOptions(config, productId)
	if err != nil {
		return
	}

	bug = zentaoHelper.PrepareBug(commConsts.WorkDir, resultDir, caseId, productId)

	err = zentaoHelper.CommitBug(bug, config)
	return
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
