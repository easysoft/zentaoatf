package zentaoService

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/service/client"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"strconv"
)

func SubmitResult() {
	//conf := configUtils.ReadCurrConfig()
	Login("conf.Url", "conf.Account", "conf.Password")

	report := testingService.GetTestTestReportForSubmit(vari.CurrScriptFile, vari.CurrResultDate)

	for idx, cs := range report.Cases {
		id := cs.Id
		idInTask := cs.IdInTask

		//uri := fmt.Sprintf("testtask-runCase-%d-%d-1.json", idInTask, id)

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

		url := "conf.Url + uri"
		_, ok := client.PostObject(url, requestObj)
		if ok {
			resultId := GetLastResult("conf.Url", idInTask, id)
			report.Cases[idx].ZentaoResultId = resultId

			json, _ := json.Marshal(report)
			testingService.SaveTestTestReportAfterSubmit(vari.CurrScriptFile, vari.CurrResultDate, string(json))

			logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("success_to_submit_result", id, resultId) + "\n")
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
