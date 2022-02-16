package zentaoUtils

import (
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

func Login(config commDomain.ProjectConf) (err error) {
	err = GetConfig(config.Url)
	if err != nil {
		logUtils.Infof(i118Utils.Sprintf("fail_to_login"))
		return
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
	bodyBytes, err = httpUtils.PostStr(url, params)
	if err != nil || (err == nil && strings.Index(string(bodyBytes), "title") > 0) { // use PostObject to login again for new system
		_, err = httpUtils.Post(url, params, true)
	}

	if err == nil {
		if commConsts.Verbose {
			logUtils.Info(i118Utils.Sprintf("success_to_login"))
		}
	} else {
		logUtils.Errorf(i118Utils.Sprintf("fail_to_login"))
	}

	return
}

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

	// check site path by calling login interface
	uri := ""
	if commConsts.RequestType == commConsts.PathInfo {
		uri = "user-login.json"
	} else {
		uri = "index.php?m=user&f=login&t=json"
	}
	url = baseUrl + uri
	bytes, err = httpUtils.Get(url)
	if err != nil {
		return
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
		err = errors.New("请先完成项目配置")
		return
	}

	err = Login(config)
	if err != nil {
		return
	}

	// $productID = 0, $branch = 0, $browseType = '', $param = 0, $storyType = 'story',
	// $orderBy = '', $recTotal = 0, $recPerPage = 20, $pageID = 1, $projectID = 0)
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("-----id_asc-0-10000-1-0")
	} else {
		params = fmt.Sprintf("orderBy=id_asc&recTotal=0&recPerPage=10000")
	}

	url := config.Url + GenApiUri("product", "browse", params)
	bytes, err := httpUtils.Get(url)

	if err != nil {
		err = errors.New("请检查项目配置")
		return
	}

	jsn, _ := simplejson.NewJson(bytes)
	productMap, err := jsn.Get("products").Map()

	for key, val := range productMap {
		id, _ := strconv.Atoi(key)
		products = append(products, serverDomain.ZentaoProduct{Id: id, Name: val.(string)})
	}

	return
}

func ListModuleByProduct(productId int, projectPath string) (modules []serverDomain.ZentaoModule, err error) {
	config := configUtils.LoadByProjectPath(projectPath)
	err = Login(config)
	if err != nil {
		return
	}

	// tree-browse-1-story.html#app=product
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-story", productId)
	} else {
		params = fmt.Sprintf("rootID=%d&viewType=story", productId)
	}

	url := config.Url + GenApiUri("tree", "browse", params)
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

	url := config.Url + GenApiUri("tree", "browse", params)
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

	url := config.Url + GenApiUri("testsuite", "browse", params)

	bytes, err := httpUtils.Get(url)
	if err != nil {
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
	err = Login(config)
	if err != nil {
		return
	}

	// $productID = 0, $branch = '', $type = 'local,totalStatus', $orderBy = 'id_asc', $recTotal = 0, $recPerPage = 20, $pageID = 1, $beginTime = 0, $endTime = 0)
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d--local,totalStatus-id_asc-0-10000-1", productId)
	} else {
		params = fmt.Sprintf("productID=%d&type=local,totalStatus&orderBy=id_asc&recTotal=0&recPerPage=10000", productId)
	}

	url := config.Url + GenApiUri("testtask", "browse", params)
	bytes, err := httpUtils.Get(url)

	if err != nil {
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
