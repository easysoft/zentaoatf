package zentaoHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	"github.com/easysoft/zentaoatf/internal/comm/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	"github.com/easysoft/zentaoatf/internal/pkg/domain"
	httpUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"path"
)

func GetConfig(baseUrl string) (err error) {
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

func Login(config commDomain.WorkspaceConf) (err error) {
	url := GenApiUrl("tokens", nil, config.Url)

	params := map[string]string{
		"account":  config.Username,
		"password": config.Password,
	}
	bodyBytes, err := httpUtils.Post(url, params)
	if err != nil {
		err = ZentaoLoginErr(err.Error())
		return
	}

	jsn, err := simplejson.NewJson(bodyBytes)
	if err != nil {
		err = ZentaoLoginErr(err.Error())

		return
	}

	if jsn == nil {
		return
	}
	mp, err := jsn.Map()
	if err != nil {
		err = ZentaoLoginErr(err.Error())

		return
	}

	val, ok := mp["token"]
	if ok {
		commConsts.SessionId = val.(string)
		if commConsts.Verbose {
			logUtils.Info(i118Utils.Sprintf("success_to_login"))
		}

	} else {
		err = ZentaoLoginErr(fmt.Sprintf("err response: %#v", string(bodyBytes)))

		return
	}

	return
}

func ListLang() (langs []serverDomain.ZentaoLang, err error) {
	for key, mp := range commConsts.LangMap {
		langs = append(langs, serverDomain.ZentaoLang{Code: key, Name: mp["name"]})
	}

	return
}

func LoadSiteProduct(currSite serverDomain.ZentaoSite, currProductId int) (
	products []serverDomain.ZentaoProduct, currProduct serverDomain.ZentaoProduct, err error) {

	if currSite.Url == "" {
		products = []serverDomain.ZentaoProduct{}
		return
	}
	config := commDomain.WorkspaceConf{
		Url:      currSite.Url,
		Username: currSite.Username,
		Password: currSite.Password,
	}

	products, err = loadProduct(config)
	if err != nil {
		products = make([]serverDomain.ZentaoProduct, 0)
		return
	}

	currIndex := 0
	for idx, item := range products {
		product := serverDomain.ZentaoProduct{Id: item.Id, Name: item.Name}

		if currProductId == product.Id {
			currIndex = idx
		}
	}

	currProduct = products[currIndex]
	//products[currIndex].Checked = true

	return
}

func ListProduct(workspacePath string) (products []serverDomain.ZentaoProduct, err error) {
	config := configHelper.LoadByWorkspacePath(workspacePath)
	return loadProduct(config)
}

func loadProduct(config commDomain.WorkspaceConf) (products []serverDomain.ZentaoProduct, err error) {
	if config.Url == "" {
		err = errors.New(i118Utils.Sprintf("pls_config_workspace"))
		return
	}

	err = Login(config)
	if err != nil {
		return
	}

	url := GenApiUrl("products", nil, config.Url)
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
	items, err := jsn.Get("products").Array()
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	products = []serverDomain.ZentaoProduct{}
	for _, item := range items {
		productMap, _ := item.(map[string]interface{})

		id, _ := productMap["id"].(json.Number).Int64()
		name, _ := productMap["name"].(string)

		products = append(products, serverDomain.ZentaoProduct{Id: int(id), Name: name})
	}

	return
}

func ListCaseModule(productId uint, site model.Site) (modules []domain.NestedItem, err error) {
	config := configHelper.LoadBySite(site)
	return LoadCaseModule(productId, config)
}
func LoadCaseModule(productId uint, config commDomain.WorkspaceConf) (modules []domain.NestedItem, err error) {
	err = Login(config)
	if err != nil {
		return
	}

	arr, err := LoadCaseModuleArr(productId, config)
	if err != nil {
		return
	}

	modules = make([]domain.NestedItem, 0)
	for _, item := range arr {
		genModuleListData(item, "", &modules)
	}

	return
}

func LoadCaseModuleArr(productId uint, config commDomain.WorkspaceConf) (arr []interface{}, err error) {
	err = Login(config)
	if err != nil {
		return
	}

	//uri := fmt.Sprintf("products/%d?fields=modules", productId) // this is product modules
	uri := fmt.Sprintf("modules?type=case&id=%d", productId)
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

	arr, err = jsn.Get("modules").Array()
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	return
}

func genModuleListData(interf interface{}, parentName string, modules *[]domain.NestedItem) {
	mp := interf.(map[string]interface{})

	idNum := mp["id"].(json.Number)
	id, _ := idNum.Int64()
	name := path.Join("/", parentName, mp["name"].(string))
	*modules = append(*modules, domain.NestedItem{Id: int(id), Name: name})

	if mp["children"] == nil {
		return
	}

	children := mp["children"].([]interface{})
	for _, child := range children {
		childMap := child.(map[string]interface{})
		genModuleListData(childMap, name, modules)
	}
}

func ListSuite(productId uint, site model.Site) (products []domain.NestedItem, err error) {
	config := configHelper.LoadBySite(site)
	return LoadSuite(productId, config)
}
func LoadSuite(productId uint, config commDomain.WorkspaceConf) (suites []domain.NestedItem, err error) {
	err = Login(config)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("products/%d/testsuites", productId)
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

	arr, err := jsn.Get("testsuites").Array()
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	suites, _ = GenPlatItems(arr)

	return
}

func ListTask(productId uint, site model.Site) (products []domain.NestedItem, err error) {
	config := configHelper.LoadBySite(site)
	return LoadTask(productId, config)
}
func LoadTask(productId uint, config commDomain.WorkspaceConf) (tasks []domain.NestedItem, err error) {
	if config.Url == "" {
		err = errors.New(i118Utils.Sprintf("pls_config_workspace"))
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
		err = ZentaoRequestErr(err.Error())
		return
	}

	jsn, err := simplejson.NewJson(bytes)
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	arr, err := jsn.Get("testtasks").Array()
	if err != nil {
		err = ZentaoRequestErr(url, commConsts.ResponseParseErr.Message)
		return
	}

	tasks, _ = GenPlatItems(arr)

	return
}

func GenPlatItems(arr []interface{}) (ret []domain.NestedItem, err error) {
	for _, iterf := range arr {
		temp := iterf.(map[string]interface{})
		id64, _ := temp["id"].(json.Number).Int64()

		item := domain.NestedItem{Id: int(id64), Name: temp["name"].(string)}
		ret = append(ret, item)
	}

	return
}
