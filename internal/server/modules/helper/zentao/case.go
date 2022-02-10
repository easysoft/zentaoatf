package zentaoUtils

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/helper/config"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/helper/script"
	"github.com/emirpasic/gods/maps"
	"sort"
	"strconv"
	"strings"
)

func GetCasesByModule(productId int, moduleId int, projectPath string) (cases []string) {
	config := configUtils.LoadByProjectPath(projectPath)
	ok := Login(config)
	if !ok {
		return
	}

	testcases := ListCaseByModule(config.Url, productId, moduleId)

	caseIdMap := map[int]string{}
	for _, tc := range testcases {
		id, _ := strconv.Atoi(tc.Id)
		caseIdMap[id] = ""
	}

	commonUtils.ChangeScriptForDebug(&projectPath)
	scriptUtils.GetScriptByIdsInDir(projectPath, caseIdMap, &cases)

	return
}

func GetCasesBySuite(productId int, suiteId int, projectPath string) (cases []string) {
	config := configUtils.LoadByProjectPath(projectPath)
	ok := Login(config)
	if !ok {
		return
	}

	testcases := ListCaseBySuite(config.Url, productId, suiteId)

	caseIdMap := map[int]string{}
	for _, tc := range testcases {
		id, _ := strconv.Atoi(tc.Id)
		caseIdMap[id] = ""
	}

	commonUtils.ChangeScriptForDebug(&projectPath)
	scriptUtils.GetScriptByIdsInDir(projectPath, caseIdMap, &cases)

	return
}

func GetCasesByTask(productId int, taskId int, projectPath string) (cases []string) {
	config := configUtils.LoadByProjectPath(projectPath)
	ok := Login(config)
	if !ok {
		return
	}

	testcases := ListCaseByTask(config.Url, productId, taskId)

	caseIdMap := map[int]string{}
	for _, tc := range testcases {
		id, _ := strconv.Atoi(tc.Id)
		caseIdMap[id] = ""
	}

	commonUtils.ChangeScriptForDebug(&projectPath)
	scriptUtils.GetScriptByIdsInDir(projectPath, caseIdMap, &cases)

	return
}

func ListCaseByProduct(baseUrl string, productId int) []commDomain.ZtfCase {
	// $productID=productId, $branch = '', $browseType = 'byModule', $param=moduleId,
	// $orderBy='id_desc', $recTotal=0, $recPerPage=10000, $pageID=1)

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d--byModule-all-id_asc-0-10000-1", productId)
	} else {
		params = fmt.Sprintf("productID=%d&branch=&browseType=byModule&param=0&orderBy=id_desc&recTotal=0&recPerPage=10000", productId)
	}

	url := baseUrl + GenApiUri("testcase", "browse", params)
	dataStr, ok := httpUtils.Get(url)

	if ok {
		var product commDomain.ZtfProduct
		json.Unmarshal(dataStr, &product)

		caseArr := make([]commDomain.ZtfCase, 0)
		for _, cs := range product.Cases {
			caseId := cs.Id

			csWithSteps := GetCaseById(baseUrl, caseId)
			stepArr := genCaseSteps(csWithSteps)
			caseArr = append(caseArr, commDomain.ZtfCase{Id: caseId, Product: cs.Product, Module: cs.Module,
				Title: cs.Title, StepArr: stepArr})
		}

		return caseArr
	}

	return nil
}

func ListCaseByModule(baseUrl string, productId, moduleId int) []commDomain.ZtfCase {
	// $productID = 0, $branch = '', $browseType = 'bymodule', $param = 0, $orderBy = 'id_asc', $recTotal = 0, $recPerPage = 20, $pageID = 1, $projectID = 0)
	// testcase-browse-1--byModule-19

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d--bymodule-%d-id_asc-0-10000-1-0", productId, moduleId)
	} else {
		params = fmt.Sprintf("productID=%d&browseType=bymodule&param=%d&orderBy=id_asc&recTotal=0&recPerPage=10000", productId, moduleId)
	}

	url := baseUrl + GenApiUri("testcase", "browse", params)
	bytes, ok := httpUtils.Get(url)

	if ok {
		var module commDomain.ZtfModule
		json.Unmarshal(bytes, &module)

		caseArr := make([]commDomain.ZtfCase, 0)
		for _, cs := range module.Cases {
			caseId := cs.Id

			csWithSteps := GetCaseById(baseUrl, caseId)
			stepArr := genCaseSteps(csWithSteps)

			caseArr = append(caseArr, commDomain.ZtfCase{Id: caseId, Product: cs.Product, Module: cs.Module,
				Title: cs.Title, StepArr: stepArr})
		}

		return caseArr
	}

	return nil
}

func ListCaseBySuite(baseUrl string, productId, suiteId int) []commDomain.ZtfCase {
	// $suiteID, $orderBy = 'id_asc', $recTotal = 0, $recPerPage = 20, $pageID = 1

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-id_asc-0-10000-1", suiteId)
	} else {
		params = fmt.Sprintf("suiteID=%d&orderBy=id_asc&recTotal=0&recPerPage=10000", suiteId)
	}

	url := baseUrl + GenApiUri("testsuite", "view", params)
	bytes, ok := httpUtils.Get(url)

	if ok {
		var suite commDomain.ZtfSuite
		json.Unmarshal(bytes, &suite)

		caseArr := make([]commDomain.ZtfCase, 0)
		for _, cs := range suite.Cases {
			caseId := cs.Id

			csWithSteps := GetCaseById(baseUrl, caseId)
			stepArr := genCaseSteps(csWithSteps)

			caseArr = append(caseArr, commDomain.ZtfCase{Id: caseId, Product: cs.Product, Module: cs.Module,
				Title: cs.Title, StepArr: stepArr})
		}

		return caseArr
	}

	return nil
}

func ListCaseByTask(baseUrl string, productId, taskId int) []commDomain.ZtfCase {
	// $taskID, $browseType = 'all', $param = 0,
	// $orderBy = 'id_asc', $recTotal = 0, $recPerPage = 20, $pageID = 1

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-all-0-id_asc-0-10000-1", taskId)
	} else {
		params = fmt.Sprintf("taskID=%d&browseType=all&param=0&orderBy=id_asc&recTotal=0&recPerPage=10000", taskId)
	}

	url := baseUrl + GenApiUri("testtask", "cases", params)
	bytes, ok := httpUtils.Get(url)

	if ok {
		var task commDomain.ZtfTask
		json.Unmarshal(bytes, &task)

		caseArr := make([]commDomain.ZtfCase, 0)
		for _, cs := range task.Runs {
			caseId := cs.Case

			csWithSteps := GetCaseById(baseUrl, caseId)
			stepArr := genCaseSteps(csWithSteps)

			caseArr = append(caseArr, commDomain.ZtfCase{Id: caseId, Product: cs.Product, Module: cs.Module,
				Title: cs.Title, StepArr: stepArr})
		}

		return caseArr
	}

	return nil
}

func genCaseSteps(csWithSteps commDomain.ZtfCase) (ret []commDomain.ZtfStep) {
	// get order keys
	keys := make([]int, 0, len(csWithSteps.Steps))
	for k := range csWithSteps.Steps {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		step := csWithSteps.Steps[key]
		ret = append(ret, step)
	}

	return
}

func GetCaseById(baseUrl string, caseId string) commDomain.ZtfCase {
	// $caseID, $version = 0, $from = 'testcase', $taskID = 0

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%s-0-testcase-0", caseId)
	} else {
		params = fmt.Sprintf("caseID=%s&version=0&$from=testcase&taskID=0", caseId)
	}

	url := baseUrl + GenApiUri("testcase", "view", params)
	bytes, ok := httpUtils.Get(url)

	if ok {
		var csw commDomain.ZtfCaseWrapper
		json.Unmarshal(bytes, &csw)

		cs := csw.Case
		return cs
	}

	return commDomain.ZtfCase{}
}

func CommitCase(caseId int, title string,
	stepMap maps.Map, stepTypeMap maps.Map, expectMap maps.Map, projectPath string) {
	config := configUtils.LoadByProjectPath(projectPath)

	ok := Login(config)
	if !ok {
		return
	}

	// $caseID, $comment = false
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-0", caseId)
	} else {
		params = fmt.Sprintf("caseID=%d&comment=0", caseId)
	}

	url := config.Url + GenApiUri("testcase", "edit", params)

	requestObj := map[string]interface{}{"title": title,
		"steps":    commonUtils.LinkedMapToMap(stepMap),
		"stepType": commonUtils.LinkedMapToMap(stepTypeMap),
		"expects":  commonUtils.LinkedMapToMap(expectMap)}

	json, _ := json.Marshal(requestObj)

	if commConsts.Verbose {
		logUtils.Infof(string(json))
	}

	_, ok = httpUtils.Post(url, requestObj, true)
	if ok {
		logUtils.Infof(i118Utils.Sprintf("success_to_commit_case", caseId) + "\n")
	}
}

func fieldMapToListOrderByInt(mp map[string]interface{}) []commDomain.BugOption {
	arr := make([]commDomain.BugOption, 0)

	keys := make([]int, 0)
	for key, _ := range mp {
		keyint, _ := strconv.Atoi(key)
		keys = append(keys, keyint)
	}

	sort.Ints(keys)

	for _, key := range keys {
		keyStr := strconv.Itoa(key)
		name := strings.TrimSpace(mp[keyStr].(string))
		if name == "" {
			name = "-"
		}

		opt := commDomain.BugOption{Id: keyStr, Name: name}
		arr = append(arr, opt)
	}

	return arr
}

func fieldMapToListOrderByStr(mp map[string]interface{}, notNull bool) []commDomain.BugOption {
	arr := make([]commDomain.BugOption, 0)

	keys := make([]string, 0)
	for key, _ := range mp {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		name := strings.TrimSpace(mp[key].(string))
		if name == "" {
			if notNull {
				continue
			}
		}

		opt := commDomain.BugOption{Id: key, Name: name}
		arr = append(arr, opt)
	}

	return arr
}

func fieldArrToListKeyStr(arr0 []interface{}, notNull bool) []commDomain.BugOption {
	arr := make([]commDomain.BugOption, 0)

	keys := make([]string, 0)
	for _, val := range arr0 {
		keys = append(keys, val.(string))
	}

	sort.Strings(keys)

	for _, val := range arr0 {
		name := val.(string)
		if name == "" {
			if notNull {
				continue
			}
		}

		opt := commDomain.BugOption{Id: val.(string), Name: name}
		arr = append(arr, opt)
	}

	return arr
}
