package zentaoUtils

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/comm/helper/analysis"
	"github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/bitly/go-simplejson"
	"github.com/fatih/color"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
)

func CommitBug(ztfBug commDomain.ZtfBug, projectPath string) (err error) {
	config := configUtils.LoadByProjectPath(projectPath)
	Login(config)

	ztfBug.Steps = strings.Replace(ztfBug.Steps, " ", "&nbsp;", -1)
	ztfBug.Steps = strings.Replace(ztfBug.Steps, "\n", "<br />", -1)

	// bug-create-1-0-caseID=1,version=3,resultID=93,runID=0,stepIdList=9_12_
	// bug-create-1-0-caseID=1,version=3,resultID=84,runID=6,stepIdList=9_12_,testtask=2,projectID=1,buildID=1
	// http://zentaopms.deeptest.com/bug-create-1-0-moduleID=0.html
	extras := fmt.Sprintf("caseID=%s,version=%s,resultID=0,runID=0,stepIdList=%s",
		ztfBug.Case, ztfBug.Version, ztfBug.StepIds)

	// $productID, $branch = '', $extras = ''
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%s-0-%s", ztfBug.Product, extras)
	} else {
		params = fmt.Sprintf("productID=%s&branch=0&$extras=%s", ztfBug.Product, extras)
	}
	//params = ""
	url := config.Url + GenApiUri("bug", "create", params)

	bug := commDomain.ZentaoBug{}
	copier.Copy(&bug, ztfBug)
	ret, ok := httpUtils.Post(url, bug, true)

	msg := ""

	if ok {
		msg = i118Utils.Sprintf("success_to_report_bug", ztfBug.Case)
	} else {
		msg = color.RedString(string(ret))
	}

	if commConsts.ComeFrom == "cmd" {
		msgView, _ := commConsts.Cui.View("reportBugMsg")
		msgView.Clear()
		if ok {
			color.New(color.FgGreen).Fprintf(msgView, msg)

			commConsts.Cui.DeleteView("submitInput")

			cancelReportBugInput, _ := commConsts.Cui.View("cancelReportBugInput")
			cancelReportBugInput.Clear()
			fmt.Fprint(cancelReportBugInput, " "+i118Utils.Sprintf("close"))
		} else {
			color.New(color.FgMagenta).Fprintf(msgView, msg)
		}
	}

	logUtils.Info(msg)

	return
}

func PrepareBug(projectPath, seq string, caseIdStr string) (bug commDomain.ZtfBug) {
	caseId, err := strconv.Atoi(caseIdStr)

	if err != nil {
		return
	}

	report, err := analysisUtils.ReadReport(projectPath, seq)
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
			Product: strconv.Itoa(cs.ProductId), Case: strconv.Itoa(cs.Id),
			Uid:   uuid.NewV4().String(),
			Steps: strings.Join(steps, "\n"), StepIds: stepIds,
			Version: "trunk", Severity: "3", Pri: "3",
			OpenedBuild: map[string]string{"0": "trunk"}, CaseVersion: "0", OldTaskID: "0",
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

func GetBugFiledOptions(req commDomain.FuncResult, projectPath string) (
	bugFields commDomain.ZentaoBugFields, err error) {

	// field options
	config := configUtils.LoadByProjectPath(projectPath)
	ok := Login(config)
	if !ok {
		return
	}

	// $productID, $projectID = 0
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-0", req.ProductId)
	} else {
		params = fmt.Sprintf("productID=%d", req.ProductId)
	}

	url := config.Url + GenApiUri("bug", "ajaxGetBugFieldOptions", params)
	bytes, ok := httpUtils.Get(url)
	if ok {
		jsonData := &simplejson.Json{}
		jsonData, err = simplejson.NewJson(bytes)

		if err != nil {
			return
		}

		mp, _ := jsonData.Get("modules").Map()
		bugFields.Modules = fieldMapToListOrderByInt(mp)

		mp, _ = jsonData.Get("categories").Map()
		bugFields.Categories = fieldMapToListOrderByStr(mp, false)

		mp, _ = jsonData.Get("versions").Map()
		bugFields.Versions = fieldMapToListOrderByStr(mp, true)

		mp, _ = jsonData.Get("severities").Map()
		bugFields.Severities = fieldMapToListOrderByInt(mp)

		arr, _ := jsonData.Get("priorities").Array()
		bugFields.Priorities = fieldArrToListKeyStr(arr, true)
	}

	return
}
