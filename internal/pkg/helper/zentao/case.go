package zentaoHelper

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"

	"github.com/fatih/color"

	"github.com/bitly/go-simplejson"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stdinUtils "github.com/easysoft/zentaoatf/pkg/lib/stdin"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"github.com/kataras/iris/v12"
)

func CommitCase(caseId int, title string, steps []commDomain.ZentaoCaseStep, script serverDomain.TestScript,
	config commDomain.WorkspaceConf, noNeedConfirm, withCode bool) (err error) {

	if serverConfig.CONFIG.AuthToken == "" {
		err = Login(config)
		if err != nil {
			return
		}
	}

	_, err = GetCaseById(config, caseId)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("/testcases/%d", caseId)
	url := GenApiUrl(uri, nil, config.Url)

	requestObj := map[string]interface{}{
		"type": "feature",
	}

	requestObj["title"] = title
	requestObj["steps"] = steps
	if withCode {
		requestObj["path"] = script.Path
		requestObj["script"] = script.Code
		requestObj["lang"] = script.Lang
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
	if !noNeedConfirm && commConsts.ExecFrom == commConsts.FromCmd {
		logUtils.ExecConsole(color.FgCyan, "\n"+i118Utils.Sprintf("case_update_confirm", caseId, title))
		stdinUtils.InputForBool(&yes, true, "want_to_continue")
	}

	if yes {
		_, err = httpUtils.Put(url, requestObj)
		if err != nil {
			err = ZentaoRequestErr(url, err.Error())
			return
		}

		logUtils.Infof(i118Utils.Sprintf("success_to_commit_case", caseId))
	}

	return
}

func CreateCase(productId int, title string, steps []commDomain.ZentaoCaseStep, script serverDomain.TestScript,
	config commDomain.WorkspaceConf) (cs commDomain.ZtfCase, err error) {

	err = Login(config)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("/products/%d/testcases", productId)
	url := GenApiUrl(uri, nil, config.Url)

	requestObj := map[string]interface{}{
		"type":   "feature",
		"title":  title,
		"steps":  steps,
		"script": script.Code,
		"lang":   script.Lang,
		"pri":    3,
	}
	if steps == nil {
		requestObj["steps"] = []commDomain.ZentaoCaseStep{}
	}

	jsonData, err := json.Marshal(requestObj)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	if commConsts.Verbose {
		logUtils.Infof(string(jsonData))
	}

	bytes, err := httpUtils.Post(url, requestObj)
	if err != nil {
		err = ZentaoRequestErr(url, err.Error())
		return
	}
	err = json.Unmarshal(bytes, &cs)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}
	logUtils.Infof(i118Utils.Sprintf("success_to_create_case", title) + "\n")

	return
}

func GetCaseById(config commDomain.WorkspaceConf, caseId int) (cs commDomain.ZtfCase, err error) {
	if serverConfig.CONFIG.AuthToken == "" {
		err = Login(config)
		if err != nil {
			return
		}
	}

	uri := fmt.Sprintf("/testcases/%d", caseId)
	url := GenApiUrl(uri, nil, config.Url)
	logUtils.Infof("req url: %s", url)
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

func LoadTestCasesInModuleTree(workspace model.Workspace, scriptIdsFromZentao map[int]string,
	productId int, suiteId, taskId int, config commDomain.WorkspaceConf) (root serverDomain.TestAsset, err error) {

	// scripts in workspace dir
	scriptsInDir, _ := scriptHelper.LoadScriptTreeByDir(workspace, scriptIdsFromZentao)
	scriptsMapInDir := map[int]*serverDomain.TestAsset{}
	genScriptMap(scriptsInDir, &scriptsMapInDir)

	// cases in zentao by module
	casesFromZentao, err := LoadTestCaseSimple(productId, 0, suiteId, taskId, config)
	caseMapByModuleFromZentao := map[int][]commDomain.ZtfCaseInModule{}
	for _, item := range casesFromZentao.Cases {
		caseMapByModuleFromZentao[item.Module] = append(caseMapByModuleFromZentao[item.Module], item)
	}

	// modules in zentao
	modules, err := LoadCaseModuleArr(uint(productId), config)

	root = serverDomain.TestAsset{
		Type:          commConsts.Workspace,
		WorkspaceId:   int(workspace.ID),
		WorkspaceType: workspace.Type,
		Path:          workspace.Path,
		Title:         fileUtils.GetDirName(workspace.Path),
		Slots:         iris.Map{"icon": "icon"},

		Checkable: true,
		IsLeaf:    false,
	}

	for _, item := range modules {
		genModuleTreeWithCases(item, scriptsMapInDir, caseMapByModuleFromZentao, &root)
	}

	return
}

func genScriptMap(asset serverDomain.TestAsset, mp *map[int]*serverDomain.TestAsset) {
	if asset.CaseId > 0 {
		(*mp)[asset.CaseId] = &asset
	}

	if asset.Children == nil {
		return
	}

	for _, child := range asset.Children {
		genScriptMap(*child, mp)
	}

	return
}

func genModuleTreeWithCases(moduleInterface interface{},
	casesMapInDir map[int]*serverDomain.TestAsset, caseMapByModuleFromZentao map[int][]commDomain.ZtfCaseInModule,
	asset *serverDomain.TestAsset) {

	moduleMap := moduleInterface.(map[string]interface{})

	idNum := moduleMap["id"].(json.Number)
	id64, _ := idNum.Int64()
	moduleId := int(id64)
	moduleName := moduleMap["name"].(string)

	// add module node
	dirNode := scriptHelper.AddDir("", moduleId, moduleName, asset)

	// add cases in module
	for _, cs := range caseMapByModuleFromZentao[moduleId] {
		caseId := cs.Id
		caseNameInZentao := cs.Title
		casePath := ""

		// case info from dir
		_, ok := casesMapInDir[caseId]
		if ok {
			casePath = casesMapInDir[caseId].Path
		}

		scriptHelper.AddScript(moduleId, caseId, casePath, caseNameInZentao, "module", true, dirNode)
	}

	if moduleMap["children"] == nil {
		return
	}

	children := moduleMap["children"].([]interface{})
	for _, child := range children {
		genModuleTreeWithCases(child, casesMapInDir, caseMapByModuleFromZentao, dirNode)
	}
}

func LoadTestCaseSimple(productId, moduleId, suiteId, taskId int,
	config commDomain.WorkspaceConf) (casesResp commDomain.ZtfRespTestCases, err error) {

	err = Login(config)
	if err != nil {
		return
	}

	if moduleId != 0 {
		casesResp, err = ListCaseByModule(config.Url, productId, moduleId)
	} else if suiteId != 0 {
		casesResp, err = ListCaseBySuite(config.Url, suiteId)
	} else if taskId != 0 {
		casesResp, err = ListCaseByTask(config.Url, taskId)
	} else if productId != 0 { // low priority at below
		casesResp, err = ListCaseByProduct(config.Url, productId)
	}

	return
}

func LoadTestCasesDetail(productId, moduleId, suiteId, taskId int,
	config commDomain.WorkspaceConf) (cases []commDomain.ZtfCase, err error) {

	casesResp, _ := LoadTestCaseSimple(productId, moduleId, suiteId, taskId, config)

	for _, cs := range casesResp.Cases {
		caseId := cs.Id

		cs, err := GetTestCaseDetail(caseId, config)
		if err != nil {
			continue
		}

		cases = append(cases, cs)
	}

	return
}

func LoadTestCasesDetailByCaseIds(caseIds []int,
	config commDomain.WorkspaceConf) (cases []commDomain.ZtfCase, err error) {

	for _, caseId := range caseIds {
		cs, err := GetTestCaseDetail(caseId, config)
		if err != nil {
			continue
		}

		cases = append(cases, cs)
	}

	return
}

func GetTestCaseDetail(caseId int, config commDomain.WorkspaceConf) (cs commDomain.ZtfCase, err error) {
	csWithSteps, err := GetCaseById(config, caseId)

	stepArr := genCaseSteps(csWithSteps)
	cs = commDomain.ZtfCase{Id: caseId, Product: csWithSteps.Product, Module: csWithSteps.Module,
		Title: csWithSteps.Title, Steps: stepArr}

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
		id, _ := strconv.ParseInt(fmt.Sprintf("%v", mp["id"]), 10, 64)

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
		id, _ := strconv.ParseInt(fmt.Sprintf("%v", mp["id"]), 10, 64)

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
	url := GenApiUrl(uri, map[string]interface{}{"limit": 10000}, config.Url)

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
		id, _ := strconv.ParseInt(fmt.Sprintf("%v", mp["case"]), 10, 64)

		caseIdMap[int(id)] = ""
	}

	return
}

func GetCasesByModuleInDir(productId, moduleId int, workspacePath, scriptDir string) (cases []string, err error) {
	config := commDomain.WorkspaceConf{}
	config = configHelper.LoadByWorkspacePath(workspacePath)

	err = Login(config)
	if err != nil {
		return
	}

	caseResp, err := ListCaseByModule(config.Url, productId, moduleId)
	if err != nil {
		return
	}

	caseIds := make([]int, 0)
	for _, tc := range caseResp.Cases {
		caseIds = append(caseIds, tc.Id)
	}

	caseIdMap := map[int]string{}
	scriptHelper.GetScriptByIdsInDir(scriptDir, &caseIdMap)

	cases = scriptHelper.GetCaseByListInMap(caseIds, caseIdMap)

	return
}

func GetCasesBySuiteInDir(productId int, suiteId int, workspacePath, scriptDir string) (cases []string, err error) {
	config := commDomain.WorkspaceConf{}
	config = configHelper.LoadByWorkspacePath(workspacePath)

	err = Login(config)
	if err != nil {
		return
	}

	caseResp, err := ListCaseBySuite(config.Url, suiteId)
	if err != nil {
		return
	}

	caseIds := make([]int, 0)
	for _, cs := range caseResp.Cases {
		caseIds = append(caseIds, cs.Id)
	}

	caseIdMap := map[int]string{}
	scriptHelper.GetScriptByIdsInDir(scriptDir, &caseIdMap)

	cases = scriptHelper.GetCaseByListInMap(caseIds, caseIdMap)

	return
}

func GetCasesByTaskInDir(productId int, taskId int, workspacePath, scriptDir string) (cases []string, err error) {
	config := commDomain.WorkspaceConf{}
	config = configHelper.LoadByWorkspacePath(workspacePath)

	err = Login(config)
	if err != nil {
		return
	}

	caseResp, err := ListCaseByTask(config.Url, taskId)
	if err != nil {
		return
	}

	caseIds := make([]int, 0)
	for _, tc := range caseResp.Cases {
		caseIds = append(caseIds, tc.Case)
	}

	caseIdMap := map[int]string{}
	scriptHelper.GetScriptByIdsInDir(scriptDir, &caseIdMap)

	cases = scriptHelper.GetCaseByListInMap(caseIds, caseIdMap)

	return
}

func ListCaseByProduct(baseUrl string, productId int) (casesResp commDomain.ZtfRespTestCases, err error) {
	uri := fmt.Sprintf("/products/%d/testcases", productId)
	url := GenApiUrl(uri, map[string]interface{}{
		"limit": 10000,
	}, baseUrl)

	dataStr, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	casesResp, err = parseCases(dataStr)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	return
}

func parseCases(bytes []byte) (casesResp commDomain.ZtfRespTestCases, err error) {
	casesResp = commDomain.ZtfRespTestCases{}
	casesResp.Cases = make([]commDomain.ZtfCaseInModule, 0)
	jsn, err := simplejson.NewJson(bytes)
	if err != nil {
		return
	}

	items, err := jsn.Get("testcases").Array()
	parseCaseItem(items, &casesResp)
	return
}

func parseCaseItem(items []interface{}, casesResp *commDomain.ZtfRespTestCases) {
	for _, item := range items {
		itemMap := item.(map[string]interface{})
		isScene, ok := itemMap["isScene"]
		if ok {
			cases, casesExist := itemMap["cases"]
			if casesExist {
				casesMap := cases.(map[string]interface{})
				for _, caseItem := range casesMap {
					casesResp.Cases = append(casesResp.Cases, convertMap2Case(caseItem))
				}
			}

			if isScene.(bool) {
				children, childrenExist := itemMap["children"]
				if childrenExist {
					childItems := []interface{}{}
					childrenMap := children.(map[string]interface{})
					for _, childItem := range childrenMap {
						childItems = append(childItems, childItem)
					}
					parseCaseItem(childItems, casesResp)
				}

				continue
			}
		}

		casesResp.Cases = append(casesResp.Cases, convertMap2Case(item))
	}
}

func convertMap2Case(caseItem interface{}) (caseStruct commDomain.ZtfCaseInModule) {
	caseItemMap := caseItem.(map[string]interface{})

	_, ok := caseItemMap["id"].(json.Number)
	if ok {
		caseId, _ := caseItemMap["id"].(json.Number).Int64()
		caseStruct.Id = int(caseId)
	} else {
		caseStruct.Id, _ = strconv.Atoi(strings.Replace(caseItemMap["id"].(string), "case_", "", 1))
	}

	product, _ := caseItemMap["product"].(json.Number).Int64()
	caseStruct.Product = int(product)
	module, _ := caseItemMap["module"].(json.Number).Int64()
	caseStruct.Module = int(module)
	if caseItemMap["case"] != nil {
		taskCaseId, _ := caseItemMap["case"].(json.Number).Int64()
		caseStruct.Case = int(taskCaseId)
	}
	caseStruct.Title = stringUtils.ToStr(caseItemMap["title"])
	caseStruct.StatusName = stringUtils.ToStr(caseItemMap["statusName"])
	caseStruct.Type = stringUtils.ToStr(caseItemMap["type"])

	caseStruct.LastRunResult = caseItemMap["lastRunResult"].(string)

	return
}

func ListCaseByModule(baseUrl string, productId, moduleId int) (casesResp commDomain.ZtfRespTestCases, err error) {
	uri := fmt.Sprintf("/products/%d/testcases", productId)
	url := GenApiUrl(uri, map[string]interface{}{
		"module": moduleId,
		"limit":  10000,
	}, baseUrl)
	spew.Dump(url)
	bytes, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}
	spew.Dump(string(bytes))
	casesResp, err = parseCases(bytes)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	return
}

func ListCaseBySuite(baseUrl string, suiteId int) (casesResp commDomain.ZtfRespTestCases, err error) {
	uri := fmt.Sprintf("/testsuites/%d", suiteId)
	url := GenApiUrl(uri, map[string]interface{}{
		"limit": 10000,
	}, baseUrl)

	bytes, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	casesResp, err = parseCases(bytes)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	return
}

func ListCaseByTask(baseUrl string, taskId int) (casesResp commDomain.ZtfRespTestCases, err error) {
	uri := fmt.Sprintf("/testtasks/%d", taskId)
	url := GenApiUrl(uri, map[string]interface{}{
		"limit": 10000,
	}, baseUrl)

	bytes, err := httpUtils.Get(url)
	if err != nil {
		err = ZentaoRequestErr(err.Error())
		return
	}

	casesResp, err = parseCases(bytes)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}
	for index, caseInfo := range casesResp.Cases {
		casesResp.Cases[index].Id = caseInfo.Case
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
