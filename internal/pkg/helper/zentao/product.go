package zentaoHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strconv"

	"github.com/bitly/go-simplejson"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/pkg/domain"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
)

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

	if len(products) == 0 {
		return
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

	url := GenApiUrl("products", map[string]interface{}{"limit": 1000}, config.Url)
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

		id, _ := strconv.ParseInt(fmt.Sprintf("%v", productMap["id"]), 10, 64)
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
		id64, _ := strconv.ParseInt(fmt.Sprintf("%v", temp["id"]), 10, 64)

		item := domain.NestedItem{Id: int(id64), Name: temp["name"].(string)}
		ret = append(ret, item)
	}

	return
}
