package zentaoHelper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"

	"github.com/bitly/go-simplejson"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	analysisHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/analysis"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stdinUtils "github.com/easysoft/zentaoatf/pkg/lib/stdin"
	"github.com/fatih/color"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
)

func CommitBug(ztfBug commDomain.ZtfBug, config commDomain.WorkspaceConf) (err error) {
	if ztfBug.Product == 0 {
		logUtils.Info(color.RedString(i118Utils.Sprintf("ignore_bug_without_product")))
		return
	}

	Login(config)

	if ztfBug.TestType == commConsts.TestUnit {
		return CommitUnitBug(&ztfBug, config)
	}
	ztfBug.Steps = strings.Replace(ztfBug.Steps, " ", "&nbsp;", -1)
	ztfBug.Steps = strings.Replace(ztfBug.Steps, "\n", "<br />", -1)

	uri := fmt.Sprintf("/products/%d/bugs", ztfBug.Product)
	url := GenApiUrl(uri, nil, config.Url)

	bug := commDomain.ZentaoBug{}
	copier.Copy(&bug, ztfBug)
	_, err = httpUtils.Post(url, bug)
	if err != nil {
		err = ZentaoRequestErr(url, i118Utils.Sprintf("fail_to_report_bug", err.Error()))
		return
	}

	logUtils.Info(color.GreenString(i118Utils.Sprintf("success_to_report_bug", ztfBug.Case)))

	return
}

func CommitUnitBug(ztfBug *commDomain.ZtfBug, config commDomain.WorkspaceConf) (err error) {
	if ztfBug.Product == 0 {
		logUtils.Info(color.RedString(i118Utils.Sprintf("ignore_bug_without_product")))
		return
	}

	if ztfBug.Title == "" {
		ztfBug.Title = stdinUtils.GetInput("\\w+", "",
			i118Utils.Sprintf("pls_enter")+" "+i118Utils.Sprintf("bug_title"))
	}

	Login(config)

	submitBugs := make([]*commDomain.ZentaoBug, 0)
	ztfBug.Steps = strings.Replace(ztfBug.Steps, " ", "&nbsp;", -1)
	ztfBug.Steps = strings.Replace(ztfBug.Steps, "\n", "<br />", -1)
	ztfBug.Title = strings.Trim(ztfBug.Title, "-")

	bug := commDomain.ZentaoBug{}
	copier.Copy(&bug, ztfBug)
	submitBugs = append(submitBugs, &bug)
	generateCaseId(submitBugs, config)
	for _, bug := range submitBugs {
		uri := fmt.Sprintf("/products/%d/bugs", bug.Product)
		url := GenApiUrl(uri, nil, config.Url)

		_, err = httpUtils.Post(url, bug)
		if err != nil {
			err = ZentaoRequestErr(url, i118Utils.Sprintf("fail_to_report_bug", err.Error()))
			return
		}

		logUtils.Info(color.GreenString(i118Utils.Sprintf("success_to_report_bug", bug.Case)))
	}

	return
}

func PrepareBug(workspacePath, seq string, caseIdStr string, productId int) (bug commDomain.ZtfBug) {
	caseId, err := strconv.Atoi(caseIdStr)

	if err != nil {
		return
	}

	report, _, err := analysisHelper.ReadReportByWorkspaceSeq(workspacePath, seq)
	if err != nil {
		return
	}

	if report.TestType == commConsts.TestFunc {
		return prepareBugForFunc(report, productId, caseId)
	}

	//unit type
	return prepareBugForUnit(report, productId, caseId)
}

func prepareBugForFunc(report commDomain.ZtfReport, productId, caseId int) (bug commDomain.ZtfBug) {
	for _, cs := range report.FuncResult {
		if cs.Id != caseId || cs.Status != commConsts.FAIL {
			continue
		}

		steps, stepIds := parseStep(cs)

		bug = commDomain.ZtfBug{
			Title:    cs.Title,
			TestType: "func",
			Case:     cs.Id,
			Product:  productId,
			Steps:    steps,
			StepIds:  stepIds,

			Uid:  uuid.NewV4().String(),
			Type: "codeerror", Severity: 3, Pri: 3, OpenedBuild: []string{"trunk"},
			CaseVersion: "0", OldTaskID: "0",
		}
	}

	return
}

func prepareBugForUnit(report commDomain.ZtfReport, productId, caseId int) (bug commDomain.ZtfBug) {
	bugs := make([]map[string]interface{}, 0)
	stepIds := ""

	for _, cs := range report.UnitResult {
		if (caseId > 0 && cs.Id != caseId) || cs.Status == "pass" {
			continue
		}
		stepIds += strconv.Itoa(cs.Id) + "_"
		steps := ""
		if cs.Failure != nil {
			steps = cs.Failure.Desc
		}

		bugMap := map[string]interface{}{
			"title":  cs.Title,
			"status": cs.Status,
			"steps":  steps,
			"caseId": cs.Cid,
			"id":     cs.Id,
		}
		bugs = append(bugs, bugMap)
		break
	}

	bugsJSon, _ := json.Marshal(bugs)

	bug = commDomain.ZtfBug{
		Title:    report.Name,
		TestType: "unit",
		Case:     caseId,
		Product:  productId,
		Steps:    string(bugsJSon),
		StepIds:  stepIds,

		Uid:  uuid.NewV4().String(),
		Type: "codeerror", Severity: 3, Pri: 3, OpenedBuild: []string{"trunk"},
		CaseVersion: "0", OldTaskID: "0",
	}

	return
}

func parseStep(cs commDomain.FuncResult) (bugSteps string, stepIds string) {
	steps := make([]string, 0)
	stepsArray := make([]map[string]interface{}, 0)

	for _, step := range cs.Steps {
		if step.Status == commConsts.FAIL {
			stepIds += step.Id + "_"
		}
		stepIds = strings.Trim(stepIds, "_")
		stepsContent := GenBugStepText(step)
		steps = append(steps, stepsContent)
		stepsArray = append(stepsArray, map[string]interface{}{
			"title":  step.Name,
			"status": step.Status,
			"steps":  stepsContent,
		})
	}

	if commConsts.ExecFrom == commConsts.FromClient {
		jsonSteps, _ := json.Marshal(stepsArray)
		bugSteps = string(jsonSteps)
		return
	}

	bugSteps = strings.Join(steps, "\n")
	return
}

func GenBugStepText(step commDomain.StepLog) string {
	stepResults := make([]string, 0)

	stepId, _ := strconv.Atoi(step.Id)
	stepTxt := fmt.Sprintf("%s%d： %s %s\n", i118Utils.Sprintf("step"), stepId, step.Name, step.Status)
	for _, checkpoint := range step.CheckPoints {
		text := fmt.Sprintf(
			"  %s：%s\n"+
				"    %s：\n"+
				"      %s\n"+
				"    %s：\n"+
				"      %s",
			i118Utils.Sprintf("checkpoint"), checkpoint.Status, i118Utils.Sprintf("expect_result"), checkpoint.Expect, i118Utils.I118Prt.Sprintf("actual_result"), checkpoint.Actual)

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
		err = ZentaoRequestErr(url, err.Error())
		return
	}

	bytes = []byte(strings.Replace(string(bytes), `"modules":["\/"]`, `"modules":{"0":"\/"}`, 1))
	err = json.Unmarshal(bytes, &bugOptionsWrapper)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	arr, ok := bugOptionsWrapper.Options.SeverityObj.([]interface{})
	if ok {
		bugFields.Severity = fieldArrToListKeyStr(arr, false)
	} else {
		mp, ok := bugOptionsWrapper.Options.SeverityObj.(iris.Map)
		if ok {
			bugFields.Severity = fieldMapToListOrderByInt(mp)
		}
	}

	bugFields.Types = fieldMapToListOrderByStr(bugOptionsWrapper.Options.Type, false)
	bugFields.Pri = fieldArrToListKeyStr(bugOptionsWrapper.Options.Pri, false)

	bugFields.Modules = fieldMapToListOrderByInt(bugOptionsWrapper.Options.Modules)
	bugFields.Build = fieldMapToListOrderByStr(bugOptionsWrapper.Options.Build, true)

	return
}

func LoadBugs(Product int, config commDomain.WorkspaceConf) (bugs []commDomain.ZentaoBug, err error) {

	Login(config)

	uri := fmt.Sprintf("/products/%d/bugs", Product)
	params := map[string]interface{}{
		"limit": 10000,
	}
	url := GenApiUrl(uri, params, config.Url)

	bytes, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(url, err.Error())
		return
	}
	jsn, err := simplejson.NewJson(bytes)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}
	items, err := jsn.Get("bugs").Array()

	for _, item := range items {
		bug, _ := item.(map[string]interface{})
		id, _ := strconv.ParseInt(fmt.Sprintf("%v", bug["id"]), 10, 64)
		caseId, _ := strconv.ParseInt(fmt.Sprintf("%v", bug["case"]), 10, 64)
		severity, _ := strconv.ParseInt(fmt.Sprintf("%v", bug["severity"]), 10, 64)
		pri, _ := strconv.ParseInt(fmt.Sprintf("%v", bug["pri"]), 10, 64)
		title, _ := bug["title"].(string)
		openedBy, _ := bug["openedBy"].(map[string]interface{})
		bugs = append(bugs, commDomain.ZentaoBug{
			Id:         int(id),
			Title:      title,
			Case:       int(caseId),
			Steps:      bug["steps"].(string),
			Type:       bug["type"].(string),
			Severity:   int(severity),
			Pri:        int(pri),
			StatusName: bug["statusName"].(string),
			OpenedDate: bug["openedDate"].(string),
			OpenedBy:   openedBy["realname"].(string),
		})
	}
	return
}

func generateCaseId(bugs []*commDomain.ZentaoBug, config commDomain.WorkspaceConf) (caseId int) {
	if len(bugs) == 0 {
		return
	}
	//查询所有case，标题相同则返回
	casesResp, _ := ListCaseByProduct(config.Url, bugs[0].Product)
	for _, bug := range bugs {
		for _, cs := range casesResp.Cases {
			if cs.Title == bug.Title {
				bug.Case = cs.Id
				break
			}
		}
		if bug.Case == 0 {
			caseInfo, _ := CreateCase(fmt.Sprintf("%v", bug.Product), bug.Title, nil, serverDomain.TestScript{}, config)
			bug.Case = caseInfo.Id
		}
	}

	return
}
