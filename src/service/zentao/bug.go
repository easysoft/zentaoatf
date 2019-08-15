package zentaoService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/service/client"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	printUtils "github.com/easysoft/zentaoatf/src/utils/print"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

func SubmitBug() {
	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	productId := conf.ProductId
	projectId := conf.ProjectId

	report := testingService.GetTestTestReportForSubmit(vari.CurrScriptFile, vari.CurrResultDate)
	for _, cs := range report.Cases {
		if cs.Id != vari.CurrCaseId {
			continue
		}

		id := cs.Id
		idInTask := cs.IdInTask
		taskId := cs.TaskId
		zentaoResultId := cs.ZentaoResultId
		title := cs.Title

		// bug-create-1-0-caseID=1,version=3,resultID=93,runID=0,stepIdList=9_12_
		// bug-create-1-0-caseID=1,version=3,resultID=84,runID=6,stepIdList=9_12_,testtask=2,projectID=1,buildID=1
		stepIds := ""

		requestObj := map[string]interface{}{"module": "0", "uid": uuid.NewV4().String(),
			"caseVersion": "0", "oldTaskID": "0", "product": productId, "project": projectId,
			"case": cs.Id, "result": cs.ZentaoResultId, "title": title,
		}

		version := map[string]interface{}{"0": "trunk"}
		requestObj["openedBuild"] = version

		if taskId != 0 {
			requestObj["testtask"] = taskId
		}

		for _, step := range cs.Steps {
			if !step.Status {
				stepIds += strconv.Itoa(step.Id) + "_"
			}

			requestObj["steps"] = testingService.GetStepHtml(step)
		}

		params := fmt.Sprintf("caseID=%d,version=0,resultID=%d,runID=%d,stepIdList=%s",
			id, zentaoResultId, idInTask, stepIds)

		if taskId != 0 {
			temp := fmt.Sprintf("testtask=%d,projectID=%d,buildID=1", taskId, projectId)
			params += temp
		}

		uri := fmt.Sprintf("bug-create-%d-0-%s.json", productId, params)
		printUtils.PrintToCmd(uri)

		reqStr, _ := json.Marshal(requestObj)
		printUtils.PrintToCmd(string(reqStr))

		url := conf.Url + uri
		_, ok := client.PostObject(url, requestObj)
		if ok {
			printUtils.PrintToCmd(
				fmt.Sprintf("success to submit a bug for case %d-%d", id, idInTask))
		}
	}
}

func GetZentaoSettings() {
	config := configUtils.ReadCurrConfig()

	entityType := config.EntityType
	entityVal := config.EntityVal

	requestObj := make(map[string]interface{})
	requestObj["entityType"] = entityType
	requestObj["entityVal"] = entityVal

	url := config.Url
	_, _ = client.PostJson(url+constant.UrlZentaoSettings, requestObj)
	printUtils.PrintToCmd(url + constant.UrlZentaoSettings)

	if err == nil {
		if pass {
			utils.PrintToCmd("success to get settings")
			//utils.ZendaoSettings = body.ZentaoSettings
		}
	} else {
		printUtils.PrintToCmd(err.Error())
	}
}
