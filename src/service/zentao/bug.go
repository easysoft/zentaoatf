package zentaoService

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
)

func PrepareBug(resultDir string, caseIdStr string) (model.Bug, string) {
	caseId, err := strconv.Atoi(caseIdStr)

	if err != nil {
		return model.Bug{}, ""
	}

	report := testingService.GetZTFTestReportForSubmit(resultDir)
	for _, cs := range report.FuncResult {
		if cs.Id != caseId {
			continue
		}

		product := cs.ProductId
		GetBugFiledOptions(product)

		title := cs.Title
		module := GetFirstNoEmptyVal(vari.ZentaoBugFileds.Modules)
		typ := GetFirstNoEmptyVal(vari.ZentaoBugFileds.Categories)
		openedBuild := map[string]string{"0": "trunk"}
		severity := GetFirstNoEmptyVal(vari.ZentaoBugFileds.Severities)
		priority := GetFirstNoEmptyVal(vari.ZentaoBugFileds.Priorities)

		caseId := cs.Id

		uid := uuid.NewV4().String()
		caseVersion := "0"
		oldTaskID := "0"

		stepIds := ""
		steps := make([]string, 0)
		for _, step := range cs.Steps {
			if !step.Status {
				stepIds += step.Id + "_"
			}

			stepsContent := testingService.GetStepContent(step)
			steps = append(steps, stepsContent)
		}

		bug := model.Bug{Title: title,
			Module: module, Type: typ, OpenedBuild: openedBuild, Severity: severity, Pri: priority,
			Product: strconv.Itoa(product), Case: strconv.Itoa(caseId),
			Steps: strings.Join(steps, "\n"),
			Uid:   uid, CaseVersion: caseVersion, OldTaskID: oldTaskID,
		}
		return bug, stepIds
	}

	return model.Bug{}, ""
}

func CommitBug() (bool, string) {
	bug := vari.CurrBug
	stepIds := vari.CurrBugStepIds

	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	productId := bug.Product
	bug.Steps = strings.Replace(bug.Steps, " ", "&nbsp;", -1)

	// bug-create-1-0-caseID=1,version=3,resultID=93,runID=0,stepIdList=9_12_
	// bug-create-1-0-caseID=1,version=3,resultID=84,runID=6,stepIdList=9_12_,testtask=2,projectID=1,buildID=1
	p := fmt.Sprintf("caseID=%s,version=0,resultID=0,runID=0,stepIdList=%s",
		bug.Case, stepIds)

	// $productID, $branch = '', $extras = ''
	params := ""
	if vari.RequestType == constant.RequestTypePathInfo {
		params = fmt.Sprintf("%s-0-%s.json", productId, p)
	} else {
		params = fmt.Sprintf("productID=%s&branch=0&$extras=%s", productId, p)
	}

	url := conf.Url + zentaoUtils.GenApiUri("bug", "create", params)

	body, ok := client.PostObject(url, bug)
	if !ok {
		return false, ""
	}

	json, err1 := simplejson.NewJson([]byte(body))
	if err1 != nil {
		return false, ""
	}

	msg, err2 := json.Get("message").String()
	if err2 != nil {
		return false, ""
	}

	if msg == "" {
		return true, i118Utils.I118Prt.Sprintf("success_to_report_bug", bug.Case)
	} else {
		return false, msg
	}
}
