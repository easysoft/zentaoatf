package zentaoService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"sort"
	"strconv"
)

func LoadTestCases(productId string, moduleId string, suiteIdStr string, taskIdStr string) []model.TestCase {
	config := configUtils.ReadCurrConfig()

	var testcases []model.TestCase

	ok := Login(config.Url, config.Account, config.Password)
	if !ok {
		return testcases
	}

	if moduleId != "" {
		testcases = ListCaseByModule(config.Url, productId, moduleId)
	} else if suiteIdStr != "" {
		testcases = ListCaseBySuite(config.Url, suiteIdStr)
	} else if taskIdStr != "" {
		testcases = ListCaseByTask(config.Url, taskIdStr)
	} else if productId != "" {
		testcases = ListCaseByProduct(config.Url, productId)
	} else {
		logUtils.PrintUsage()
	}

	return testcases
}

func ListCaseByProduct(baseUrl string, productId string) []model.TestCase {
	params := [][]string{{"productID", productId}}
	url := baseUrl + zentaoUtils.GenSuperApiUri("testcase", "getModuleCases", params)
	dataStr, ok := client.Get(url, nil)

	if ok {
		var caseMap map[string]model.TestCase
		json.Unmarshal([]byte(dataStr), &caseMap)

		keys := make([]int, 0)
		for _, cs := range caseMap {
			i, _ := strconv.Atoi(cs.Id)
			keys = append(keys, i)
		}
		sort.Ints(keys)

		caseArr := make([]model.TestCase, 0)
		for _, id := range keys {
			idStr := strconv.Itoa(id)

			tc := caseMap[idStr]
			csWithSteps := GetCaseById(baseUrl, idStr)
			caseArr = append(caseArr, model.TestCase{Id: idStr, Product: tc.Product, Title: tc.Title, StepArr: csWithSteps.StepArr})
		}

		return caseArr
	}

	return nil
}

func ListCaseByModule(baseUrl string, productId string, moduleId string) []model.TestCase {
	params := fmt.Sprintf("%s--byModule-%s-id_asc-0-10000-1", productId, moduleId)

	url := baseUrl + zentaoUtils.GenApiUri("testcase", "browse", params)
	dataStr, ok := client.Get(url, nil)

	if ok {
		var module model.Module
		json.Unmarshal([]byte(dataStr), &module)

		caseArr := make([]model.TestCase, 0)
		for _, cs := range module.Cases {
			caseId := cs.Id

			csWithSteps := GetCaseById(baseUrl, caseId)
			caseArr = append(caseArr, model.TestCase{Id: caseId, Product: cs.Product,
				Title: cs.Title, StepArr: csWithSteps.StepArr})
		}

		return caseArr
	}

	return nil
}

func ListCaseBySuite(baseUrl string, suiteId string) []model.TestCase {
	params := fmt.Sprintf("%s-id_asc-0-10000-1", suiteId)

	url := baseUrl + zentaoUtils.GenApiUri("testsuite", "view", params)
	dataStr, ok := client.Get(url, nil)

	if ok {
		var task model.TestSuite
		json.Unmarshal([]byte(dataStr), &task)

		caseArr := make([]model.TestCase, 0)
		for _, cs := range task.Cases {
			caseId := cs.Id

			csWithSteps := GetCaseById(baseUrl, caseId)
			caseArr = append(caseArr, model.TestCase{Id: caseId, Product: cs.Product,
				Title: cs.Title, StepArr: csWithSteps.StepArr})
		}

		return caseArr
	}

	return nil
}

func ListCaseByTask(baseUrl string, suiteId string) []model.TestCase {
	params := fmt.Sprintf("%s-all-0-id_asc-0-10000-1", suiteId)

	url := baseUrl + zentaoUtils.GenApiUri("testtask", "cases", params)
	dataStr, ok := client.Get(url, nil)

	if ok {
		var task model.TestTask
		json.Unmarshal([]byte(dataStr), &task)

		caseArr := make([]model.TestCase, 0)
		for _, cs := range task.Runs {
			caseId := cs.Case

			csWithSteps := GetCaseById(baseUrl, caseId)
			caseArr = append(caseArr, model.TestCase{Id: caseId, Product: cs.ProductId,
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

func GetCaseIdsBySuite(suiteId string, idMap *map[int]string) {
	config := configUtils.ReadCurrConfig()

	ok := Login(config.Url, config.Account, config.Password)
	if !ok {
		return
	}

	testcases := ListCaseBySuite(config.Url, suiteId)

	for _, tc := range testcases {
		id, _ := strconv.Atoi(tc.Id)
		(*idMap)[id] = ""
	}
}

func GetCaseIdsByTask(taskId string, idMap *map[int]string) {
	config := configUtils.ReadCurrConfig()

	ok := Login(config.Url, config.Account, config.Password)
	if !ok {
		return
	}

	testcases := ListCaseByTask(config.Url, taskId)

	for _, tc := range testcases {
		id, _ := strconv.Atoi(tc.Id)
		(*idMap)[id] = ""
	}
}
