package zentaoUtils

import (
	"errors"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/helper/config"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/bitly/go-simplejson"
	"strconv"
	"strings"
)

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
			logUtils.Info(i118Utils.Sprintf("success_to_login"))
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

func ListLang() (langs []serverDomain.ZentaoLang, err error) {
	for key, _ := range langUtils.LangMap {
		langs = append(langs, serverDomain.ZentaoLang{Code: key, Name: key})
	}

	return
}

func ListProduct(projectPath string) (products []serverDomain.ZentaoProduct, err error) {
	config := configUtils.LoadByProjectPath(projectPath)
	if config.Url == "" {
		err = errors.New("请先完成项目配置")
		return
	}

	Login(config)

	// $productID = 0, $branch = 0, $browseType = '', $param = 0, $storyType = 'story',
	// $orderBy = '', $recTotal = 0, $recPerPage = 20, $pageID = 1, $projectID = 0)
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("-----id_asc-0-10000-1-0")
	} else {
		params = fmt.Sprintf("orderBy=id_asc&recTotal=0&recPerPage=10000")
	}

	url := config.Url + GenApiUri("product", "browse", params)
	bytes, ok := httpUtils.Get(url)

	if !ok {
		err = errors.New("请检查项目配置")
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

	// $productID = 0, $orderBy = 'id_asc', $recTotal = 0, $recPerPage = 20, $pageID = 1
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-id_asc-0-10000-1", productId)
	} else {
		params = fmt.Sprintf("productID=%d&orderBy=id_asc&recTotal=0&recPerPage=10000", productId)
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

	// $productID = 0, $branch = '', $type = 'local,totalStatus', $orderBy = 'id_asc', $recTotal = 0, $recPerPage = 20, $pageID = 1, $beginTime = 0, $endTime = 0)
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d--local,totalStatus-id_asc-0-10000-1", productId)
	} else {
		params = fmt.Sprintf("productID=%d&type=local,totalStatus&orderBy=id_asc&recTotal=0&recPerPage=10000", productId)
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

func GetBugFiledOptions(req commDomain.FuncResult, projectPath string) (
	bugFields commDomain.ZentaoBugFields, err error) {

	// field options
	config := configUtils.LoadByProjectPath(projectPath)
	ok := Login(config)
	if !ok {
		return
	}

	// $productID, $projectID = 0
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-0", req.ProductId)
	} else {
		params = fmt.Sprintf("productID=%d", req.ProductId)
	}

	url := config.Url + GenApiUri("bug", "ajaxGetBugFieldOptions", params)
	bytes, ok := httpUtils.Get(url)
	if ok {
		jsonData := &simplejson.Json{}
		jsonData, err = simplejson.NewJson(bytes)

		if err != nil {
			return
		}

		mp, _ := jsonData.Get("modules").Map()
		bugFields.Modules = fieldMapToListOrderByInt(mp)

		mp, _ = jsonData.Get("categories").Map()
		bugFields.Categories = fieldMapToListOrderByStr(mp, false)

		mp, _ = jsonData.Get("versions").Map()
		bugFields.Versions = fieldMapToListOrderByStr(mp, true)

		mp, _ = jsonData.Get("severities").Map()
		bugFields.Severities = fieldMapToListOrderByInt(mp)

		arr, _ := jsonData.Get("priorities").Array()
		bugFields.Priorities = fieldArrToListKeyStr(arr, true)
	}

	return
}

func GetStepText(step commDomain.StepLog) string {
	stepResults := make([]string, 0)

	stepTxt := fmt.Sprintf("步骤%s： %s %s\n", step.Id, step.Name, step.Status)

	for _, checkpoint := range step.CheckPoints {
		text := fmt.Sprintf(
			"  检查点：%s\n"+
				"    期待结果：\n"+
				"      %s\n"+
				"    实际结果：\n"+
				"      %s",
			checkpoint.Status, checkpoint.Expect, checkpoint.Actual)

		stepResults = append(stepResults, text)
	}

	return stepTxt + strings.Join(stepResults, "\n") + "\n"
}
