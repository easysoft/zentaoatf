package zentaoService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
)

func GenBug() (model.Bug, string, string) {
	conf := configUtils.ReadCurrConfig()
	productId := conf.ProductId
	projectId := conf.ProjectId

	report := testingService.GetTestTestReportForSubmit(vari.CurrScriptFile, vari.CurrResultDate)
	for _, cs := range report.Cases {
		if cs.Id != vari.CurrCaseId {
			continue
		}

		title := cs.Title
		module := "0"
		typ := "install"
		openedBuild := map[string]interface{}{"0": "trunk"}
		severity := "1"
		priority := "1"

		product := productId
		project := projectId
		caseId := cs.Id
		Result := cs.ZentaoResultId
		taskId := cs.TaskId

		uid := uuid.NewV4().String()
		caseVersion := ""
		oldTaskID := "0"

		idInTask := strconv.Itoa(cs.IdInTask)

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
			Product: strconv.Itoa(product), Project: strconv.Itoa(project), Case: strconv.Itoa(caseId),
			Result: strconv.Itoa(Result), Testtask: strconv.Itoa(taskId), Steps: strings.Join(steps, "<br/>"),
			Uid: uid, CaseVersion: caseVersion, OldTaskID: oldTaskID,
		}
		return bug, idInTask, stepIds

	}

	return model.Bug{}, "", ""
}

func SubmitBug(bug model.Bug, idInTask string, stepIds string) {
	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	productId := bug.Product
	projectId := bug.Project

	// bug-create-1-0-caseID=1,version=3,resultID=93,runID=0,stepIdList=9_12_
	// bug-create-1-0-caseID=1,version=3,resultID=84,runID=6,stepIdList=9_12_,testtask=2,projectID=1,buildID=1
	params := fmt.Sprintf("caseID=%s,version=0,resultID=%s,runID=%s,stepIdList=%s",
		bug.Case, bug.Result, idInTask, stepIds)

	if bug.Testtask != "" {
		temp := fmt.Sprintf("testtask=%s,projectID=%s,buildID=1", bug.Testtask, projectId)
		params += temp
	}

	uri := fmt.Sprintf("bug-create-%s-0-%s.json", productId, params)
	logUtils.PrintToCmd(uri)

	reqStr, _ := json.Marshal(bug)
	logUtils.PrintToCmd(string(reqStr))

	url := conf.Url + uri
	_, ok := client.PostObject(url, bug)
	if ok {
		logUtils.PrintToCmd(
			fmt.Sprintf("success to submit a bug for case %s-%s", bug.Case, idInTask))
	}
}
