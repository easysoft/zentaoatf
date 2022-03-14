package zentaoHelper

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/comm/helper/analysis"
	configHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
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
	config := configHelper.LoadByProjectPath(projectPath)
	Login(config)

	ztfBug.Steps = strings.Replace(ztfBug.Steps, " ", "&nbsp;", -1)
	ztfBug.Steps = strings.Replace(ztfBug.Steps, "\n", "<br />", -1)

	uri := fmt.Sprintf("/products/%d/bugs", ztfBug.Product)
	url := GenApiUrl(uri, nil, config.Url)

	bug := commDomain.ZentaoBug{}
	copier.Copy(&bug, ztfBug)
	_, err = httpUtils.Post(url, bug)

	msg := ""

	if err == nil {
		msg = i118Utils.Sprintf("success_to_report_bug", ztfBug.Case)
	} else {
		msg = color.RedString("commit bug failed, error: %s.", err.Error())
	}

	logUtils.Info(msg)

	return
}

func PrepareBug(projectPath, seq string, caseIdStr string) (bug commDomain.ZtfBug) {
	caseId, err := strconv.Atoi(caseIdStr)

	if err != nil {
		return
	}

	report, err := analysisHelper.ReadReportByProjectSeq(projectPath, seq)
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
			Product: cs.ProductId, Case: cs.Id,
			Uid:   uuid.NewV4().String(),
			Steps: strings.Join(steps, "\n"), StepIds: stepIds,
			Version: "trunk", Severity: 3, Pri: 3,
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
	config := configHelper.LoadByProjectPath(projectPath)
	err = Login(config)
	if err != nil {
		return
	}

	// $productID, $projectID = 0
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-0", req.ProductId)
	} else {
		params = fmt.Sprintf("productID=%d", req.ProductId)
	}

	url := config.Url + GenApiUriOld("bug", "ajaxGetBugFieldOptions", params)
	bytes, err := httpUtils.Get(url)
	if err == nil {
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
