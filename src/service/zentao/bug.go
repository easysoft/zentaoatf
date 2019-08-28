package zentaoService

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
)

func GenBug(resultDir string, caseIdStr string) (model.Bug, string) {
	caseId, err := strconv.Atoi(caseIdStr)

	if err != nil {
		return model.Bug{}, ""
	}

	report := testingService.GetTestTestReportForSubmit(resultDir)
	for _, cs := range report.Cases {
		if cs.Id != caseId {
			continue
		}

		title := cs.Title
		module := GetFirstNoEmptyVal(vari.ZentaoBugFileds.Modules)
		typ := GetFirstNoEmptyVal(vari.ZentaoBugFileds.Categories)
		openedBuild := map[string]string{"0": "trunk"}
		severity := GetFirstNoEmptyVal(vari.ZentaoBugFileds.Severities)
		priority := GetFirstNoEmptyVal(vari.ZentaoBugFileds.Priorities)

		product := cs.ProductId
		caseId := cs.Id

		uid := uuid.NewV4().String()
		caseVersion := "0"
		oldTaskID := "0"

		stepIds := ""
		steps := make([]string, 0)
		for _, step := range cs.Steps {
			if !step.Status {
				stepIds += strconv.Itoa(step.Id) + "_"
			}

			stepsContent := testingService.GetStepContent(step)
			steps = append(steps, stepsContent)
		}

		bug := model.Bug{Title: title,
			Module: module, Type: typ, OpenedBuild: openedBuild, Severity: severity, Pri: priority,
			Product: strconv.Itoa(product), Project: "0", Case: strconv.Itoa(caseId),
			Testtask: "0", Steps: strings.Join(steps, "<br/>"),
			Uid: uid, CaseVersion: caseVersion, OldTaskID: oldTaskID,
		}
		return bug, stepIds

	}

	return model.Bug{}, ""
}

func CommitBug(bug model.Bug, stepIds string) bool {
	// TODO: open cui

	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	productId := bug.Product
	projectId := bug.Project

	// bug-create-1-0-caseID=1,version=3,resultID=93,runID=0,stepIdList=9_12_
	// bug-create-1-0-caseID=1,version=3,resultID=84,runID=6,stepIdList=9_12_,testtask=2,projectID=1,buildID=1
	params := fmt.Sprintf("caseID=%s,version=0,resultID=%s,runID=0,stepIdList=%s",
		bug.Case, bug.Result, stepIds)

	bug.Steps = strings.Replace(bug.Steps, " ", "&nbsp;", -1)
	if bug.Testtask != "" {
		temp := fmt.Sprintf("testtask=%s,projectID=%s,buildID=1", bug.Testtask, projectId)
		params += temp
	}

	uri := fmt.Sprintf("bug-create-%s-0-%s.json", productId, params)

	url := uri
	body, ok := client.PostObject(url, bug)

	json, _ := simplejson.NewJson([]byte(body))
	msg, _ := json.Get("message").String()
	if ok && msg == "" {
		logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("success_to_report_bug", bug.Case) + "\n")

		return true
	} else {
		return false
	}
}
