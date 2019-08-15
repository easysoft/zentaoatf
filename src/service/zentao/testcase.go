package zentaoService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"sort"
)

func ListCaseByProduct(baseUrl string, productId string) []model.TestCase {
	modules := ListCaseModule(baseUrl, productId)
	_ = modules

	params := [][]string{{"productID", productId}}
	url := baseUrl + zentaoUtils.GenSuperApiUri("testcase", "getModuleCases", params)
	dataStr, ok := client.Get(url, nil)

	if ok {
		var caseMap map[string]model.TestCase
		json.Unmarshal([]byte(dataStr), &caseMap)

		caseArr := make([]model.TestCase, 0)
		for _, cs := range caseMap {
			id := cs.Id
			csWithSteps := GetCaseById(baseUrl, id)
			caseArr = append(caseArr, model.TestCase{Id: id, TaskId: "0", Title: cs.Title, StepArr: csWithSteps.StepArr})
		}

		return caseArr
	}

	return nil
}

func ListCaseByTask(baseUrl string, taskId string) []model.TestCase {
	params := fmt.Sprintf("%s-all-0-id_asc-0-10000-1", taskId)

	url := baseUrl + zentaoUtils.GenApiUri("testtask", "cases", params)
	dataStr, ok := client.Get(url, nil)

	if ok {
		var task model.TestTask
		json.Unmarshal([]byte(dataStr), &task)

		caseArr := make([]model.TestCase, 0)
		for _, cs := range task.Runs {
			caseId := cs.Case
			caseInTaskId := cs.Id

			csWithSteps := GetCaseById(baseUrl, caseId)
			caseArr = append(caseArr, model.TestCase{Id: caseId, IdInTask: caseInTaskId, TaskId: taskId,
				Title: cs.Title, StepArr: csWithSteps.StepArr})
		}

		return caseArr
	}

	return nil
}

func GetCaseById(baseUrl string, caseId string) model.TestCase {
	params := [][]string{{"caseID", caseId}}
	url := baseUrl + zentaoUtils.GenSuperApiUri("testcase", "getById", params)
	dataStr, ok := client.PostStr(url, nil)

	if ok {
		var tc model.TestCase
		json.Unmarshal([]byte(dataStr), &tc)

		var keys []int
		for key := range tc.Steps {
			keys = append(keys, key)
		}
		sort.Ints(keys)

		stepArr := make([]model.TestStep, 0)
		for _, key := range keys {
			step := tc.Steps[key]
			stepArr = append(stepArr, model.TestStep{Id: step.Id, Desc: step.Desc, Expect: step.Expect,
				Type: step.Type, Parent: step.Parent})
		}
		tc.StepArr = stepArr

		return tc
	}

	return model.TestCase{}
}
