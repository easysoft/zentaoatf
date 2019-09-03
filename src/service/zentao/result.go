package zentaoService

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/service/client"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"strconv"
)

func CommitResult(resultDir string) {
	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	report := testingService.GetTestTestReportForSubmit(resultDir)

	for _, cs := range report.Cases {
		id := cs.Id

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

		uri := fmt.Sprintf("testtask-runCase-0-%d-1.json", id)
		url := conf.Url + uri

		_, ok := client.PostObject(url, requestObj)
		if ok {
			resultId := GetLastResult(conf.Url, id)
			logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("success_to_commit_result", id, resultId)+"\n", -1)
		}
	}
}

func GetLastResult(baseUrl string, caseId int) int {
	params := fmt.Sprintf("0-%d-1.json", caseId)

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
