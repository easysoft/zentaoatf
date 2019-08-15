package zentaoService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/service/client"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	printUtils "github.com/easysoft/zentaoatf/src/utils/print"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

func SubmitBug(assert string, date string, caseId int, caseIdInTask int) {
	conf := configUtils.ReadCurrConfig()
	productId := conf.ProductId
	projectId := conf.ProjectId

	report := testingService.GetTestTestReportForSubmit(assert, date)
	for _, cs := range report.Cases {
		id := cs.Id
		idInTask := cs.IdInTask
		taskId := cs.TaskId
		zentaoResultId := cs.ZentaoResultId
		title := cs.Title

		if caseId != id || caseIdInTask != idInTask {
			continue
		}

		// bug-create-1-0-caseID=1,version=3,resultID=93,runID=0,stepIdList=9_12_.html
		// bug-create-1-0-caseID=1,version=3,resultID=84,runID=6,stepIdList=9_12_,testtask=2,projectID=1,buildID=1.html
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
		println(uri)

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
