package zentaoService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/service/client"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	printUtils "github.com/easysoft/zentaoatf/src/utils/print"
	"strconv"
)

func SubmitBug(assert string, date string, caseId int, caseIdInTask int) {
	conf := configUtils.ReadCurrConfig()
	productId := conf.ProductId
	projectId := conf.ProjectId

	report := testingService.GetTestTestReportForSubmit(assert, date)
	for idx, cs := range report.Cases {
		id := cs.Id
		idInTask := cs.IdInTask
		taskId := cs.TaskId
		zentaoResultId := cs.ZentaoResultId

		if caseId != id || caseIdInTask != idInTask {
			continue
		}

		// bug-create-1-0-caseID=1,version=3,resultID=93,runID=0,stepIdList=9_12_.html
		// bug-create-1-0-caseID=1,version=3,resultID=84,runID=6,stepIdList=9_12_,testtask=2,projectID=1,buildID=1.html
		stepList := ""
		requestObj := map[string]string{"case": strconv.Itoa(id), "version": "0"}

		for _, step := range cs.Steps {
			if !step.Status {
				stepList += strconv.Itoa(step.Id) + "_"
			}

			var stepStatus string
			if step.Status {
				stepStatus = constant.PASS.String()
			} else {
				stepStatus = constant.FAIL.String()
			}

			stepResults := ""
			for _, checkpoint := range step.CheckPoints {
				stepResults += checkpoint.Actual // strconv.FormatBool(checkpoint.Status) + ": " + checkpoint.Actual
			}

			requestObj["steps["+strconv.Itoa(step.Id)+"]"] = stepStatus
			requestObj["reals["+strconv.Itoa(step.Id)+"]"] = stepResults
		}

		params := fmt.Sprintf("caseID=%d,version=0,resultID=%d,runID=%d,stepIdList=%s",
			id, zentaoResultId, idInTask, stepList)

		if taskId != 0 {
			temp := fmt.Sprintf("testtask=%d,projectID=%d,buildID=1", taskId, projectId)
			params += temp
		}

		uri := fmt.Sprintf("bug-create-%d-0-%s.json", productId, params)

		reqStr, _ := json.Marshal(requestObj)
		printUtils.PrintToCmd(string(reqStr))

		url := conf.Url + uri
		_, ok := client.PostStr(url, requestObj)
		if ok {
			resultId := GetLastResult(conf.Url, idInTask, id)
			report.Cases[idx].ZentaoResultId = resultId

			json, _ := json.Marshal(report)
			testingService.SaveTestTestReportAfterSubmit(assert, date, string(json))

			printUtils.PrintToCmd(
				fmt.Sprintf("success to submit the results for case %d, resultId is %d", id, resultId))
		}
	}
}
