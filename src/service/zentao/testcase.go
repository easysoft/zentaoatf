package zentaoService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	stdinUtils "github.com/easysoft/zentaoatf/src/utils/stdin"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/emirpasic/gods/maps"
	"sort"
	"strconv"
	"strings"
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
			caseArr = append(caseArr, model.TestCase{Id: idStr, Product: tc.Product, Module: tc.Module,
				Title: tc.Title, StepArr: csWithSteps.StepArr})
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
			caseArr = append(caseArr, model.TestCase{Id: caseId, Product: cs.Product, Module: cs.Module,
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
			caseArr = append(caseArr, model.TestCase{Id: caseId, Product: cs.Product, Module: cs.Module,
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
			caseArr = append(caseArr, model.TestCase{Id: caseId, Product: cs.Product, Module: cs.Module,
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

		for _, key := range keys {
			stepTo := tc.Steps[key]
			testStep := model.TestStep{Id: stepTo.Id, Desc: stepTo.Desc, Expect: stepTo.Expect,
				Type: stepTo.Type, Parent: stepTo.Parent}

			tc.StepArr = append(tc.StepArr, testStep)
		}

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

func CommitCase(caseId int, title string, stepMap maps.Map, stepTypeMap maps.Map, expectMap maps.Map) {
	config := configUtils.ReadCurrConfig()

	ok := Login(config.Url, config.Account, config.Password)
	if !ok {
		return
	}

	uri := fmt.Sprintf("testcase-edit-%d.json", caseId)
	requestObj := map[string]interface{}{"title": title,
		"steps":    commonUtils.LinkedMapToMap(stepMap),
		"stepType": commonUtils.LinkedMapToMap(stepTypeMap),
		"expects":  commonUtils.LinkedMapToMap(expectMap)}

	var yes bool
	logUtils.PrintToWithColor("\n"+i118Utils.I118Prt.Sprintf("case_update_confirm", caseId, title), -1)
	stdinUtils.InputForBool(&yes, true, "want_to_continue")

	if yes {
		url := config.Url + uri
		_, ok = client.PostObject(url, requestObj)

		if ok {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("success_to_commit_case", caseId) + "\n")
		}
	}
}

func IsMutiLine(step model.TestStep) bool {
	if strings.Index(step.Desc, "\n") > -1 || strings.Index(step.Expect, "\n") > -1 {
		return true
	}

	return false
}

func GetCaseContent(stepObj model.TestStep, numb string, independentFile bool, mutiLine bool) []string {
	lines := make([]string, 0)

	step := strings.TrimSpace(stepObj.Desc)
	expect := strings.TrimSpace(stepObj.Expect)

	if independentFile {
		if mutiLine {
			expect = ">>"
		} else {
			expect = ""
		}
	}

	if mutiLine {
		lines = append(lines, fmt.Sprintf("  [%s. steps] \n%s \n  [%s. expects] \n%s",
			numb, addPrefixSpace(step, 4), numb, addPrefixSpace(expect, 4)))
	} else {
		if strings.TrimSpace(stepObj.Expect) == "" {
			lines = append(lines, fmt.Sprintf("  %s", step))
		} else {
			lines = append(lines, fmt.Sprintf("  %s >> %s", step, expect))
		}
	}

	return lines
}

func addPrefixSpace(str string, numb int) string {
	arr := strings.Split(str, "\r\n")

	ret := make([]string, 0)
	for _, line := range arr {
		line = fmt.Sprintf("%*s", numb, " ") + line

		ret = append(ret, line)
	}

	return strings.Join(ret, "\n")
}
