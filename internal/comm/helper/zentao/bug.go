package zentaoHelper

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	"github.com/easysoft/zentaoatf/internal/comm/domain"
	"github.com/easysoft/zentaoatf/internal/comm/helper/analysis"
	httpUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/http"
	"github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	"github.com/fatih/color"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
)

func CommitBug(ztfBug commDomain.ZtfBug, config commDomain.WorkspaceConf) (err error) {
	if ztfBug.Product == 0 {
		logUtils.Info(color.RedString(i118Utils.Sprintf("ignore_bug_without_product")))
		return
	}

	Login(config)

	ztfBug.Steps = strings.Replace(ztfBug.Steps, " ", "&nbsp;", -1)
	ztfBug.Steps = strings.Replace(ztfBug.Steps, "\n", "<br />", -1)

	uri := fmt.Sprintf("/products/%d/bugs", ztfBug.Product)
	url := GenApiUrl(uri, nil, config.Url)

	bug := commDomain.ZentaoBug{}
	copier.Copy(&bug, ztfBug)
	_, err = httpUtils.Post(url, bug)
	if err != nil {
		err = ZentaoRequestErr(i118Utils.Sprintf("fail_to_report_bug", err.Error()))
		return
	}

	logUtils.Info(i118Utils.Sprintf("success_to_report_bug", ztfBug.Case))

	return
}

func PrepareBug(workspacePath, seq string, caseIdStr string, productId int) (bug commDomain.ZtfBug) {
	caseId, err := strconv.Atoi(caseIdStr)

	if err != nil {
		return
	}

	report, err := analysisHelper.ReadReportByWorkspaceSeq(workspacePath, seq)
	if err != nil {
		return
	}

	for _, cs := range report.FuncResult {
		if cs.Id != caseId {
			continue
		}

		steps := make([]string, 0)
		stepIds := ""
		for _, step := range cs.Steps {
			if step.Status == commConsts.FAIL {
				stepIds += step.Id + "_"
			}

			stepsContent := GenBugStepText(step)
			steps = append(steps, stepsContent)
		}

		bug = commDomain.ZtfBug{
			Title:   cs.Title,
			Case:    cs.Id,
			Product: productId,
			Steps:   strings.Join(steps, "\n"), StepIds: stepIds,

			Uid:  uuid.NewV4().String(),
			Type: "codeerror", Severity: 3, Pri: 3, OpenedBuild: []string{"trunk"},
			CaseVersion: "0", OldTaskID: "0",
		}
		return
	}

	return
}

func GenBugStepText(step commDomain.StepLog) string {
	stepResults := make([]string, 0)

	stepTxt := fmt.Sprintf("步骤%s： %s %s\n", step.Id, step.Name, step.Status)

	for _, checkpoint := range step.CheckPoints {
		text := fmt.Sprintf(
			"  检查点：%s\n"+
				"    期待结果：\n"+
				"      %s\n"+
				"    实际结果：\n"+
				"      %s",
			checkpoint.Status, checkpoint.Expect, checkpoint.Actual)

		stepResults = append(stepResults, text)
	}

	return stepTxt + strings.Join(stepResults, "\n") + "\n"
}

func GetBugFiledOptions(config commDomain.WorkspaceConf, productId int) (
	bugFields commDomain.ZentaoBugFields, err error) {

	err = Login(config)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("/options/bug?product=%d", productId)
	url := GenApiUrl(uri, nil, config.Url)

	bytes, err := httpUtils.Get(url)
	bugOptionsWrapper := commDomain.BugOptionsWrapper{}
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	err = json.Unmarshal(bytes, &bugOptionsWrapper)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	bugFields.Types = fieldMapToListOrderByStr(bugOptionsWrapper.Options.Type, false)
	bugFields.Pri = fieldArrToListKeyStr(bugOptionsWrapper.Options.Pri, false)
	bugFields.Severity = fieldMapToListOrderByInt(bugOptionsWrapper.Options.Severity)
	bugFields.Modules = fieldMapToListOrderByInt(bugOptionsWrapper.Options.Modules)
	bugFields.Build = fieldMapToListOrderByStr(bugOptionsWrapper.Options.Build, true)

	return
}
