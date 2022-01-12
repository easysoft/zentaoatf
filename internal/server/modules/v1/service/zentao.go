package service

import (
	"errors"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/zentao"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"github.com/bitly/go-simplejson"
	"strconv"
	"strings"
)

type ZentaoService struct {
	ProjectRepo *repo.ProjectRepo `inject:""`
}

func NewZentaoService() *ZentaoService {
	return &ZentaoService{}
}

func (s *ZentaoService) ListLang() (langs []serverDomain.ZentaoLang, err error) {
	for key, _ := range langUtils.LangMap {
		langs = append(langs, serverDomain.ZentaoLang{Code: key, Name: key})
	}

	return
}

func (s *ZentaoService) ListProduct() (products []serverDomain.ZentaoProduct, err error) {
	s.Login(commConsts.ProjectConfig)

	url := commConsts.ProjectConfig.Url + zentaoUtils.GenApiUri("product", "all", "")
	url = s.AddToken(url)
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

func (s *ZentaoService) ListModuleByProduct(productId int) (modules []serverDomain.ZentaoModule, err error) {
	// tree-browse-1-story.html#app=product

	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d-story", productId)
	} else {
		params = fmt.Sprintf("rootID=%d&viewType=story", productId)
	}

	url := commConsts.ProjectConfig.Url + zentaoUtils.GenApiUri("tree", "browse", params)
	url = s.AddToken(url) + "#app=product"

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
		s.GenModuleData(mp, &modules)
	}

	return
}

func (s *ZentaoService) GenModuleData(mp map[string]interface{}, modules *[]serverDomain.ZentaoModule) {
	mpLevel := mp["level"].(int)

	idStr := mp["id"].(string)
	id, _ := strconv.Atoi(idStr)
	name := strings.Repeat("-", mpLevel) + mp["name"].(string)
	*modules = append(*modules, serverDomain.ZentaoModule{Id: id, Name: name})

	if mp["children"] == nil {
		return
	}

	children := mp["children"].([]interface{})
	for _, child := range children {
		childMap := child.(map[string]interface{})
		childMap["level"] = mp["level"].(int) + 1
		s.GenModuleData(childMap, modules)
	}
}

func (s *ZentaoService) ListSuiteByProduct(productId int) (suites []serverDomain.ZentaoSuite, err error) {
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d", productId)
	} else {
		params = fmt.Sprintf("productID=%d", productId)
	}

	url := commConsts.ProjectConfig.Url + zentaoUtils.GenApiUri("testsuite", "browse", params)
	url = s.AddToken(url)

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

func (s *ZentaoService) ListTaskByProduct(productId int) (tasks []serverDomain.ZentaoTask, err error) {
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%d", productId)
	} else {
		params = fmt.Sprintf("productID=%d", productId)
	}

	url := commConsts.ProjectConfig.Url + zentaoUtils.GenApiUri("testtask", "browse", params)
	url = s.AddToken(url)
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

func (s *ZentaoService) Login(config domain.ProjectConfig) bool {
	ok := s.GetConfig(config.Url)
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
	url = s.AddToken(url)

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

func (s *ZentaoService) GetConfig(baseUrl string) bool {
	//if commConsts.RequestType != "" {
	//	return true
	//}

	url := baseUrl + "?mode=getconfig"
	str, ok := httpUtils.Get(url)
	if !ok {
		return false
	}

	json, _ := simplejson.NewJson([]byte(str))
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
	str, ok = httpUtils.Get(url)
	if !ok {
		return false
	}

	return true
}

func (s *ZentaoService) AddToken(url string) (ret string) {
	if commConsts.RequestType == commConsts.PathInfo {
		ret = url + "?" + commConsts.SessionVar + "=" + commConsts.SessionId
	} else {
		ret = url + "&" + commConsts.SessionVar + "=" + commConsts.SessionId
	}
	ret = ret + "&XDEBUG_SESSION_START=PHPSTORM"

	return
}
