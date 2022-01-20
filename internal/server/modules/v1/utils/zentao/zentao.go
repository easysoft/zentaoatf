package zentaoUtils

import (
	"encoding/json"
	"errors"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/config"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/script"
	"github.com/bitly/go-simplejson"
	"sort"
	"strconv"
	"strings"
)

func ListLang() (langs []serverDomain.ZentaoLang, err error) {
	for key, _ := range langUtils.LangMap {
		langs = append(langs, serverDomain.ZentaoLang{Code: key, Name: key})
	}

	return
}

func ListProduct(projectPath string) (products []serverDomain.ZentaoProduct, err error) {
	config := configUtils.LoadByProjectPath(projectPath)
	Login(config)

	// $productID = 0, $branch = 0, $browseType = '', $param = 0, $storyType = 'story',
	// $orderBy = '', $recTotal = 0, $recPerPage = 20, $pageID = 1, $projectID = 0)
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("0-0----id_asc-0-10000-1-0")
	} else {
		params = fmt.Sprintf("orderBy=id_desc&recTotal=0&recPerPage=10000")
	}

	url := config.Url + GenApiUri("product", "all", params)
	bytes, ok := httpUtils.Get(url)

	if !ok {
		err = errors.New("product-all fail")
		return
	}

	jsn, _ := simplejson.NewJson(bytes)
	productMap, _ := jsn.Get("products").Map()

	for key, val := range productMap {
		id, _ := strconv.Atoi(key)
		products = append(products, serverDomain.ZentaoProduct{Id: id, Name: val.(string)})
	}

	return
}

func ListModuleByProduct(productId int, projectPath string) (modules []serverDomain.ZentaoModule, err error) {
	config := configUtils.LoadByProjectPath(projectPath)
	Login(config)
	// tree-browse-1-story.html#app=product

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-story", productId)
	} else {
		params = fmt.Sprintf("rootID=%d&viewType=story", productId)
	}

	url := config.Url + GenApiUri("tree", "browse", params)
	url += "#app=product"

	bytes, ok := httpUtils.Get(url)
	if !ok {
		err = errors.New("tree-browse-story fail")
		return
	}

	jsn, _ := simplejson.NewJson(bytes)
	arr, _ := jsn.Get("tree").Array()
	for _, item := range arr {
		mp := item.(map[string]interface{})
		mp["level"] = 0
		GenModuleData(mp, &modules)
	}

	return
}

func ListModuleForCase(productId int, projectPath string) (modules []serverDomain.ZentaoModule, err error) {
	config := configUtils.LoadByProjectPath(projectPath)
	Login(config)

	// tree-browse-1-case-0-0-qa.html
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-case-0-0-qa", productId)
	} else {
		params = fmt.Sprintf("rootID=%d&viewType=case&from=qa", productId)
	}

	url := config.Url + GenApiUri("tree", "browse", params)
	url += "#app=product"

	bytes, ok := httpUtils.Get(url)
	if !ok {
		err = errors.New("tree-browse-story fail")
		return
	}

	jsn, _ := simplejson.NewJson(bytes)
	arr, _ := jsn.Get("tree").Array()
	for _, item := range arr {
		mp := item.(map[string]interface{})
		mp["level"] = 0
		GenModuleData(mp, &modules)
	}

	return
}

func GenModuleData(mp map[string]interface{}, modules *[]serverDomain.ZentaoModule) {
	mpLevel := mp["level"].(int)

	idStr := mp["id"].(string)
	id, _ := strconv.Atoi(idStr)
	name := strings.Repeat("&nbsp;", mpLevel*3) + mp["name"].(string)
	*modules = append(*modules, serverDomain.ZentaoModule{Id: id, Name: name})

	if mp["children"] == nil {
		return
	}

	children := mp["children"].([]interface{})
	for _, child := range children {
		childMap := child.(map[string]interface{})
		childMap["level"] = mp["level"].(int) + 1
		GenModuleData(childMap, modules)
	}
}

func ListSuiteByProduct(productId int, projectPath string) (suites []serverDomain.ZentaoSuite, err error) {
	config := configUtils.LoadByProjectPath(projectPath)
	Login(config)

	// $productID = 0, $orderBy = 'id_desc', $recTotal = 0, $recPerPage = 20, $pageID = 1
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-id_asc-0-10000-1", productId)
	} else {
		params = fmt.Sprintf("productID=%d&orderBy=id_desc&recTotal=0&recPerPage=10000", productId)
	}

	url := config.Url + GenApiUri("testsuite", "browse", params)

	bytes, ok := httpUtils.Get(url)
	if !ok {
		err = errors.New("testsuite-browse fail")
		return
	}

	jsn, _ := simplejson.NewJson(bytes)
	mp, _ := jsn.Get("suites").Map()
	for _, item := range mp {
		mp := item.(map[string]interface{})

		idStr := mp["id"].(string)
		id, _ := strconv.Atoi(idStr)
		name := mp["name"].(string)

		suites = append(suites, serverDomain.ZentaoSuite{Id: id, Name: name})
	}

	return
}

func ListTaskByProduct(productId int, projectPath string) (tasks []serverDomain.ZentaoTask, err error) {
	config := configUtils.LoadByProjectPath(projectPath)
	Login(config)

	// $productID = 0, $branch = '', $type = 'local,totalStatus', $orderBy = 'id_desc', $recTotal = 0, $recPerPage = 20, $pageID = 1, $beginTime = 0, $endTime = 0)
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d--local,totalStatus-id_asc-0-10000-1", productId)
	} else {
		params = fmt.Sprintf("productID=%d&type=local,totalStatus&orderBy=id_desc&recTotal=0&recPerPage=10000", productId)
	}

	url := config.Url + GenApiUri("testtask", "browse", params)
	bytes, ok := httpUtils.Get(url)

	if !ok {
		err = errors.New("testsuite-browse fail")
		return
	}

	jsn, _ := simplejson.NewJson(bytes)
	mp, _ := jsn.Get("tasks").Map()
	for _, item := range mp {
		mp := item.(map[string]interface{})

		idStr := mp["id"].(string)
		id, _ := strconv.Atoi(idStr)
		name := mp["name"].(string)

		tasks = append(tasks, serverDomain.ZentaoTask{Id: id, Name: name})
	}

	return
}

func Login(config commDomain.ProjectConf) bool {
	ok := GetConfig(config.Url)
	if !ok {
		logUtils.Infof(i118Utils.Sprintf("fail_to_login"))
		return false
	}

	uri := ""
	if commConsts.RequestType == commConsts.PathInfo {
		uri = "user-login.json"
	} else {
		uri = "index.php?m=user&f=login&t=json"
	}
	url := config.Url + uri

	params := make(map[string]string)
	params["account"] = config.Username
	params["password"] = config.Password

	var bodyBytes []byte
	bodyBytes, ok = httpUtils.PostStr(url, params)
	if !ok || (ok && strings.Index(string(bodyBytes), "title") > 0) { // use PostObject to login again for new system
		_, ok = httpUtils.Post(url, params, true)
	}

	if ok {
		if commConsts.Verbose {
			logUtils.Infof(i118Utils.Sprintf("success_to_login"), bodyBytes)
		}
	} else {
		logUtils.Errorf(i118Utils.Sprintf("fail_to_login"))
	}

	return ok
}

func GetConfig(baseUrl string) bool {
	//if commConsts.RequestType != "" {
	//	return true
	//}

	url := baseUrl + "?mode=getconfig"
	bytes, ok := httpUtils.Get(url)
	if !ok {
		return false
	}

	json, _ := simplejson.NewJson(bytes)
	commConsts.ZenTaoVersion, _ = json.Get("version").String()
	commConsts.SessionId, _ = json.Get("sessionID").String()
	commConsts.SessionVar, _ = json.Get("sessionVar").String()
	requestType, _ := json.Get("requestType").String()
	commConsts.RequestType = commConsts.ZentaoRequestType(requestType)
	commConsts.RequestFix, _ = json.Get("requestFix").String()

	// check site path by calling login interface
	uri := ""
	if commConsts.RequestType == commConsts.PathInfo {
		uri = "user-login.json"
	} else {
		uri = "index.php?m=user&f=login&t=json"
	}
	url = baseUrl + uri
	bytes, ok = httpUtils.Get(url)
	if !ok {
		return false
	}

	return true
}

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

func ListCaseByModule(baseUrl string, productId, moduleId int) []commDomain.ZtfCase {
	// $productID = 0, $branch = '', $browseType = 'bymodule', $param = 0, $orderBy = 'id_desc', $recTotal = 0, $recPerPage = 20, $pageID = 1, $projectID = 0)
	// testcase-browse-1--byModule-19

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d--bymodule-%d-id_asc-0-10000-1-0", productId, moduleId)
	} else {
		params = fmt.Sprintf("productID=%d&browseType=bymodule&param=%d&orderBy=id_desc&recTotal=0&recPerPage=10000", productId, moduleId)
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
	// $suiteID, $orderBy = 'id_desc', $recTotal = 0, $recPerPage = 20, $pageID = 1

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-id_asc-0-10000-1", suiteId)
	} else {
		params = fmt.Sprintf("suiteID=%d&orderBy=id_desc&recTotal=0&recPerPage=10000", suiteId)
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
	// $orderBy = 'id_desc', $recTotal = 0, $recPerPage = 20, $pageID = 1

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-bymodule-0-id_asc-0-10000-1", taskId)
	} else {
		params = fmt.Sprintf("taskID=%d&browseType=all&param=0&orderBy=id_desc&recTotal=0&recPerPage=10000", taskId)
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
