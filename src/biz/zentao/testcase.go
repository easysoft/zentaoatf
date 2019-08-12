package zentao

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
)

func ListCaseByProduct(baseUrl string, productId string) []model.TestCase {
	params := map[string]string{"productID": productId}
	url := baseUrl + utils.GenSuperApiUri("testcase", "getModuleCases", params)
	bodyStr, err := http.Post(url, nil)

	if err == nil {
		var bodyJson model.ZentaoResponse
		json.Unmarshal(bodyStr, &bodyJson)

		if bodyJson.Status == "success" {
			dataStr := bodyJson.Data

			var caseMap map[string]model.TestCase
			json.Unmarshal([]byte(dataStr), &caseMap)

			caseArr := make([]model.TestCase, 0)
			for id, val := range caseMap {
				csWithSteps := GetCaseById(baseUrl, id)
				caseArr = append(caseArr, model.TestCase{Id: id, Title: val.Title, StepArr: csWithSteps.StepArr})
			}

			return caseArr
		}
	}

	return nil
}

func ListCaseByTask(baseUrl string, taskId string) []model.TestCase {

	return nil
}

func GetCaseById(baseUrl string, caseId string) model.TestCase {
	params := map[string]string{"caseID": caseId}
	url := baseUrl + utils.GenSuperApiUri("testcase", "getById", params)
	bodyStr, err := http.Post(url, nil)

	if err == nil {
		var bodyJson model.ZentaoResponse
		json.Unmarshal(bodyStr, &bodyJson)

		if bodyJson.Status == "success" {
			dataStr := bodyJson.Data

			var tc model.TestCase
			json.Unmarshal([]byte(dataStr), &tc)

			stepArr := make([]model.TestStep, 0)

			for _, step := range tc.Steps {
				stepArr = append(stepArr, model.TestStep{Id: step.Id, Desc: step.Desc, Expect: step.Expect,
					Type: step.Type, Parent: step.Parent})
			}
			tc.StepArr = stepArr

			return tc
		}
	}

	return model.TestCase{}
}
