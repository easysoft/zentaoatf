package zentaoHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/bitly/go-simplejson"
	"path"
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

func Login(config commDomain.WorkspaceConf) (err error) {
	url := GenApiUrl("tokens", nil, config.Url)

	params := map[string]string{
		"account":  config.Username,
		"password": config.Password,
	}
	bodyBytes, err := httpUtils.Post(url, params)
	if err != nil {
		logUtils.Info(i118Utils.Sprintf("fail_to_login"))
		return
	}

	if commConsts.Verbose {
		logUtils.Info(i118Utils.Sprintf("success_to_login") + string(bodyBytes))
	}

	jsn, _ := simplejson.NewJson(bodyBytes)
	if jsn == nil {
		return
	}
	mp, _ := jsn.Map()

	val, ok := mp["token"]
	if ok {
		commConsts.SessionId = val.(string)
	}

	return
}

func ListLang() (langs []serverDomain.ZentaoLang, err error) {
	for key, _ := range commConsts.LangMap {
		langs = append(langs, serverDomain.ZentaoLang{Code: key, Name: key})
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
		return
	}

	var first serverDomain.ZentaoProduct
	for idx, product := range products {
		product := serverDomain.ZentaoProduct{Id: product.Id, Name: product.Name}

		if currProductId == product.Id {
			currProduct = product
		}

		if idx == 0 {
			first = product
		}
	}

	if currProduct.Id == 0 { // not found, use the first one
		currProduct = first
	}

	return
}

func ListProduct(workspacePath string) (products []serverDomain.ZentaoProduct, err error) {
	config := configUtils.LoadByWorkspacePath(workspacePath)
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

	products = []serverDomain.ZentaoProduct{}
	for _, item := range items {
		productMap, _ := item.(map[string]interface{})

		id, _ := productMap["id"].(json.Number).Int64()
		name, _ := productMap["name"].(string)

		products = append(products, serverDomain.ZentaoProduct{Id: int(id), Name: name})
	}

	return
}

func ListModule(productId uint, site model.Site) (modules []domain.NestedItem, err error) {
	config := configUtils.LoadBySite(site)
	return LoadModule(productId, config)
}
func LoadModule(productId uint, config commDomain.WorkspaceConf) (modules []domain.NestedItem, err error) {
	err = Login(config)
	if err != nil {
		return
	}

	if config.Url == "" {
		err = errors.New(i118Utils.Sprintf("pls_config_workspace"))
		return
	}

	err = Login(config)
	if err != nil {
		return
	}

	uri := fmt.Sprintf("products/%d?fields=modules", productId)
	url := GenApiUrl(uri, nil, config.Url)

	bytes, err := httpUtils.Get(url)
	if err != nil {
		return
	}

	jsn, _ := simplejson.NewJson(bytes)
	arr, _ := jsn.Get("modules").Array()

	modules = make([]domain.NestedItem, 0)
	for _, item := range arr {
		genModuleData(item, "", &modules)
	}

	return
}

func genModuleData(interf interface{}, parentName string, modules *[]domain.NestedItem) {
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
		genModuleData(childMap, name, modules)
	}
}

func ListSuite(productId uint, site model.Site) (products []domain.NestedItem, err error) {
	config := configUtils.LoadBySite(site)
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
		return
	}

	jsn, _ := simplejson.NewJson(bytes)
	arr, _ := jsn.Get("testsuites").Array()

	suites, _ = GenPlatItems(arr)

	return
}

func ListTask(productId uint, site model.Site) (products []domain.NestedItem, err error) {
	config := configUtils.LoadBySite(site)
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
		return
	}

	jsn, _ := simplejson.NewJson(bytes)
	arr, _ := jsn.Get("testtasks").Array()

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
