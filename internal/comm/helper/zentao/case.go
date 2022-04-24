package zentaoHelper

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	"github.com/easysoft/zentaoatf/internal/comm/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	scriptHelper "github.com/easysoft/zentaoatf/internal/comm/helper/script"
	httpUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	stdinUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/stdin"
	"sort"
	"strconv"
	"strings"
)

func CommitCase(caseId int, title string,
	steps []commDomain.ZentaoCaseStep, config commDomain.WorkspaceConf) (err error) {

	err = Login(config)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("/testcases/%d", caseId)
	url := GenApiUrl(uri, nil, config.Url)

	requestObj := map[string]interface{}{
		"type":  "feature",
		"title": title,
		"steps": steps,
	}

	json, err := json.Marshal(requestObj)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	if commConsts.Verbose {
		logUtils.Infof(string(json))
	}

	yes := true
	if commConsts.ExecFrom == commConsts.FromCmd {
		logUtils.ExecConsole(1, "\n"+i118Utils.Sprintf("case_update_confirm", caseId, title))
		stdinUtils.InputForBool(&yes, true, "want_to_continue")
	}

	if yes {
		_, err = httpUtils.Put(url, requestObj)
		if err != nil {
			err = ZentaoRequestErr(url, err.Error())
			return
		}

		logUtils.Infof(i118Utils.Sprintf("success_to_commit_case", caseId) + "\n")
	}

	return
}

func GetCaseById(baseUrl string, caseId int) (cs commDomain.ZtfCase) {
	uri := fmt.Sprintf("/testcases/%d", caseId)
	url := GenApiUrl(uri, nil, baseUrl)

	bytes, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	err = json.Unmarshal(bytes, &cs)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	return
}

func LoadTestCases(productId, moduleId, suiteId, taskId int,
	config commDomain.WorkspaceConf) (cases []commDomain.ZtfCase, err error) {

	err = Login(config)
	if err != nil {
		return
	}

	if moduleId != 0 {
		cases, err = ListCaseByModule(config.Url, productId, moduleId)
	} else if suiteId != 0 {
		cases, err = ListCaseBySuite(config.Url, suiteId)
	} else if taskId != 0 {
		cases, err = ListCaseByTask(config.Url, taskId)
	} else if productId != 0 {
		cases, err = ListCaseByProduct(config.Url, productId)
	}

	return
}

func GetCaseIdsInZentaoModule(productId, moduleId uint, config commDomain.WorkspaceConf) (
	caseIdMap map[int]string, err error) {

	err = Login(config)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("/products/%d/testcases?module=%d", productId, moduleId)
	url := GenApiUrl(uri, nil, config.Url)

	bytes, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	jsn, err := simplejson.NewJson(bytes)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}
	items, err := jsn.Get("testcases").Array()
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	caseIdMap = map[int]string{}
	for _, item := range items {
		mp, _ := item.(map[string]interface{})
		id, _ := mp["id"].(json.Number).Int64()

		caseIdMap[int(id)] = ""
	}

	return
}

func GetCaseIdsInZentaoSuite(productId uint, suiteId int, config commDomain.WorkspaceConf) (
	caseIdMap map[int]string, err error) {

	err = Login(config)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("/testsuites/%d", suiteId)
	url := GenApiUrl(uri, nil, config.Url)

	bytes, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	jsn, err := simplejson.NewJson(bytes)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}
	items, err := jsn.Get("testcases").Array()
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	caseIdMap = map[int]string{}
	for _, item := range items {
		mp, _ := item.(map[string]interface{})
		id, _ := mp["id"].(json.Number).Int64()

		caseIdMap[int(id)] = ""
	}

	return
}

func GetCaseIdsInZentaoTask(productId uint, taskId int, config commDomain.WorkspaceConf) (
	caseIdMap map[int]string, err error) {

	err = Login(config)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("/testtasks/%d", taskId)
	url := GenApiUrl(uri, nil, config.Url)

	bytes, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	jsn, err := simplejson.NewJson(bytes)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}
	items, err := jsn.Get("testcases").Array()
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	caseIdMap = map[int]string{}
	for _, item := range items {
		mp, _ := item.(map[string]interface{})
		id, _ := mp["case"].(json.Number).Int64()

		caseIdMap[int(id)] = ""
	}

	return
}

func GetCasesByModuleInDir(productId, moduleId uint, workspacePath, scriptDir string) (cases []string, err error) {
	config := commDomain.WorkspaceConf{}
	config = configHelper.LoadByWorkspacePath(workspacePath)

	err = Login(config)
	if err != nil {
		return
	}

	zentaoCaseIdMap, err := GetCaseIdsInZentaoModule(productId, moduleId, config)
	if err != nil {
		return
	}

	scriptHelper.GetScriptByIdsInDir(scriptDir, zentaoCaseIdMap, &cases)

	return
}

func GetCasesBySuiteInDir(productId int, suiteId int, workspacePath, scriptDir string) (cases []string, err error) {
	config := commDomain.WorkspaceConf{}
	config = configHelper.LoadByWorkspacePath(workspacePath)

	err = Login(config)
	if err != nil {
		return
	}

	testcases, err := ListCaseBySuite(config.Url, suiteId)
	if err != nil {
		return
	}

	caseIdMap := map[int]string{}
	for _, tc := range testcases {
		caseIdMap[tc.Id] = ""
	}

	//commonUtils.ChangeScriptForDebug(&workspacePath)
	scriptHelper.GetScriptByIdsInDir(scriptDir, caseIdMap, &cases)

	return
}

func GetCasesByTaskInDir(productId int, taskId int, workspacePath, scriptDir string) (cases []string, err error) {
	config := commDomain.WorkspaceConf{}
	config = configHelper.LoadByWorkspacePath(workspacePath)

	err = Login(config)
	if err != nil {
		return
	}

	testcases, err := ListCaseByTask(config.Url, taskId)
	if err != nil {
		return
	}

	caseIdMap := map[int]string{}
	for _, tc := range testcases {
		caseIdMap[tc.Id] = ""
	}

	//commonUtils.ChangeScriptForDebug(&workspacePath)
	scriptHelper.GetScriptByIdsInDir(scriptDir, caseIdMap, &cases)

	return
}

func ListCaseByProduct(baseUrl string, productId int) (caseArr []commDomain.ZtfCase, err error) {
	uri := fmt.Sprintf("/products/%d/testcases", productId)
	url := GenApiUrl(uri, nil, baseUrl)

	dataStr, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	var cases commDomain.ZtfRespTestCases
	err = json.Unmarshal(dataStr, &cases)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	for _, cs := range cases.Cases {
		caseId := cs.Id

		csWithSteps := GetCaseById(baseUrl, caseId)
		stepArr := genCaseSteps(csWithSteps)
		caseArr = append(caseArr, commDomain.ZtfCase{Id: caseId, Product: cs.Product, Module: cs.Module,
			Title: cs.Title, Steps: stepArr})
	}

	return
}

func ListCaseByModule(baseUrl string, productId, moduleId int) (caseArr []commDomain.ZtfCase, err error) {
	uri := fmt.Sprintf("/products/%d/testcases?module=%d", productId, moduleId)
	url := GenApiUrl(uri, nil, baseUrl)

	bytes, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	var module commDomain.ZtfRespTestCases
	err = json.Unmarshal(bytes, &module)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	for _, cs := range module.Cases {
		caseId := cs.Id

		csWithSteps := GetCaseById(baseUrl, caseId)
		stepArr := genCaseSteps(csWithSteps)

		caseArr = append(caseArr, commDomain.ZtfCase{Id: caseId, Product: cs.Product, Module: cs.Module,
			Title: cs.Title, Steps: stepArr})
	}

	return
}

func ListCaseBySuite(baseUrl string, suiteId int) (caseArr []commDomain.ZtfCase, err error) {
	uri := fmt.Sprintf("/testsuites/%d", suiteId)
	url := GenApiUrl(uri, nil, baseUrl)

	bytes, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	var suite commDomain.ZtfRespTestCases
	err = json.Unmarshal(bytes, &suite)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	for _, cs := range suite.Cases {
		caseId := cs.Id

		csWithSteps := GetCaseById(baseUrl, caseId)
		stepArr := genCaseSteps(csWithSteps)

		caseArr = append(caseArr, commDomain.ZtfCase{Id: caseId, Product: cs.Product, Module: cs.Module,
			Title: cs.Title, Steps: stepArr})
	}

	return
}

func ListCaseByTask(baseUrl string, taskId int) (caseArr []commDomain.ZtfCase, err error) {
	uri := fmt.Sprintf("/testtasks/%d", taskId)
	url := GenApiUrl(uri, nil, baseUrl)

	bytes, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	var task commDomain.ZtfRespTestCases
	err = json.Unmarshal(bytes, &task)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	for _, cs := range task.Cases {
		caseId := cs.Case

		csWithSteps := GetCaseById(baseUrl, caseId)
		stepArr := genCaseSteps(csWithSteps)

		caseArr = append(caseArr, commDomain.ZtfCase{Id: caseId, Product: cs.Product, Module: cs.Module,
			Title: cs.Title, Steps: stepArr})
	}

	return
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

		opt := commDomain.BugOption{Code: keyStr, Name: name}
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
		if name == "" && notNull {
			continue
		}

		opt := commDomain.BugOption{Code: key, Name: name}
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

		opt := commDomain.BugOption{Code: val.(string), Name: name}
		arr = append(arr, opt)
	}

	return arr
}
