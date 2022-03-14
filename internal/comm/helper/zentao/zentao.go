package zentaoHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/bitly/go-simplejson"
	"strconv"
	"strings"
)

func GetConfig(baseUrl string) (err error) {
	//if commConsts.RequestType != "" {
	//	return true
	//}

	url := baseUrl + "?mode=getconfig"
	bytes, err := httpUtils.Get(url)
	if err != nil {
		return
	}

	json, _ := simplejson.NewJson(bytes)
	commConsts.ZenTaoVersion, _ = json.Get("version").String()
	commConsts.SessionId, _ = json.Get("sessionID").String()
	commConsts.SessionVar, _ = json.Get("sessionVar").String()
	requestType, _ := json.Get("requestType").String()
	commConsts.RequestType = requestType
	commConsts.RequestFix, _ = json.Get("requestFix").String()

	return
}

func Login(config commDomain.ProjectConf) (err error) {
	url := GenApiUrl("tokens", nil, config.Url)

	params := map[string]string{
		"account":  config.Username,
		"password": config.Password,
	}
	bodyBytes, err := httpUtils.Post(url, params)
	if err != nil {
		logUtils.Errorf(i118Utils.Sprintf("fail_to_login"))
	}

	if commConsts.Verbose {
		logUtils.Info(i118Utils.Sprintf("success_to_login") + string(bodyBytes))
	}

	jsn, _ := simplejson.NewJson(bodyBytes)
	mp, _ := jsn.Map()

	val, ok := mp["token"]
	if ok {
		commConsts.SessionId = val.(string)
	}

	return
}

func ListLang() (langs []serverDomain.ZentaoLang, err error) {
	for key, _ := range langUtils.LangMap {
		langs = append(langs, serverDomain.ZentaoLang{Code: key, Name: key})
	}

	return
}

func ListProduct(projectPath string) (products []serverDomain.ZentaoProduct, err error) {
	config := configUtils.LoadByProjectPath(projectPath)
	if config.Url == "" {
		err = errors.New(i118Utils.Sprintf("pls_config_project"))
		return
	}

	err = Login(config)
	if err != nil {
		return
	}

	url := GenApiUrl("products", nil, config.Url)
	bytes, err := httpUtils.Get(url)

	if err != nil {
		return
	}

	jsn, err := simplejson.NewJson(bytes)
	if err != nil {
		return
	}
	items, err := jsn.Get("products").Array()
	if err != nil {
		return
	}

	for _, item := range items {
		productMap, _ := item.(map[string]interface{})

		id, _ := productMap["id"].(json.Number).Int64()
		name, _ := productMap["name"].(string)

		products = append(products, serverDomain.ZentaoProduct{Id: int(id), Name: name})
	}

	return
}

func ListModuleForCase(productId int, projectPath string) (modules []serverDomain.ZentaoModule, err error) {
	config := configUtils.LoadByProjectPath(projectPath)
	err = Login(config)
	if err != nil {
		return
	}

	// tree-browse-1-case-0-0-qa.html
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-case-0-0-qa", productId)
	} else {
		params = fmt.Sprintf("rootID=%d&viewType=case&from=qa", productId)
	}

	url := config.Url + GenApiUriOld("tree", "browse", params)
	url += "#app=product"

	bytes, err := httpUtils.Get(url)
	if err != nil {
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
	err = Login(config)
	if err != nil {
		return
	}

	// $productID = 0, $orderBy = 'id_asc', $recTotal = 0, $recPerPage = 20, $pageID = 1
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-id_asc-0-10000-1", productId)
	} else {
		params = fmt.Sprintf("productID=%d&orderBy=id_asc&recTotal=0&recPerPage=10000", productId)
	}

	url := config.Url + GenApiUriOld("testsuite", "browse", params)

	bytes, err := httpUtils.Get(url)
	if err != nil {
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
	if config.Url == "" {
		err = errors.New(i118Utils.Sprintf("pls_config_project"))
		return
	}

	err = Login(config)
	if err != nil {
		return
	}

	params := map[string]interface{}{
		"product": productId,
		"limit":   10000,
	}
	url := GenApiUrl("testtasks", params, config.Url)
	bytes, err := httpUtils.Get(url)
	if err != nil {
		return
	}

	jsn, err := simplejson.NewJson(bytes)
	if err != nil {
		return
	}
	items, err := jsn.Get("testtasks").Array()
	if err != nil {
		return
	}

	for _, item := range items {
		taskMap, _ := item.(map[string]interface{})

		id, _ := taskMap["id"].(json.Number).Int64()
		name, _ := taskMap["name"].(string)

		tasks = append(tasks, serverDomain.ZentaoTask{Id: int(id), Name: name})
	}

	return
}
