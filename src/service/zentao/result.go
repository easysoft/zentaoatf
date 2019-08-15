package zentaoService

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/service/client"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	printUtils "github.com/easysoft/zentaoatf/src/utils/print"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"strconv"
)

func SubmitResult(assert string, date string) {
	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	report := testingService.GetTestTestReportForSubmit(assert, date)

	for idx, cs := range report.Cases {
		id := cs.Id
		idInTask := cs.IdInTask

		var uri string
		uri = fmt.Sprintf("testtask-runCase-%d-%d-1.json", idInTask, id)

		requestObj := map[string]interface{}{"case": strconv.Itoa(id), "version": "0"}

		stepMap := map[string]string{}
		realMap := map[string]string{}
		for _, step := range cs.Steps {
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
			stepMap[strconv.Itoa(step.Id)] = stepStatus
			realMap[strconv.Itoa(step.Id)] = stepResults

			requestObj["steps"] = stepMap
			requestObj["reals"] = realMap
		}

		url := conf.Url + uri
		_, ok := client.PostObject(url, requestObj)
		if ok {
			resultId := GetLastResult(conf.Url, idInTask, id)
			report.Cases[idx].ZentaoResultId = resultId

			json, _ := json.Marshal(report)
			testingService.SaveTestTestReportAfterSubmit(assert, date, string(json))

			printUtils.PrintToCmd(
				fmt.Sprintf("success to submit the results for case %d, resultId is %d", id, resultId))

			SubmitBug(assert, date, id, idInTask)
		}
	}
}

func GetLastResult(baseUrl string, caseInTaskId int, caseId int) int {
	params := fmt.Sprintf("%d-%d-1.json", caseInTaskId, caseId)

	url := baseUrl + zentaoUtils.GenApiUri("testtask", "results", params)
	dataStr, ok := client.Get(url, nil)

	resultId := -1
	if ok {
		jsonData, err := simplejson.NewJson([]byte(dataStr))
		if err == nil {
			results, _ := jsonData.Get("results").Map()

			for key, _ := range results {
				numb, _ := strconv.Atoi(key)

				if resultId < numb {
					resultId = numb
				}
			}
		}
	}

	return resultId
}
