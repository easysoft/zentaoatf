package zentaoService

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/fatih/color"
	"sort"
	"strconv"
	"strings"
)

func GetBugFiledOptions(productId int) {
	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	// $productID, $projectID = 0
	params := ""
	if vari.RequestType == constant.RequestTypePathInfo {
		params = fmt.Sprintf("%d-0", productId)
	} else {
		params = fmt.Sprintf("productID=%d", productId)
	}

	url := conf.Url + zentaoUtils.GenApiUri("bug", "ajaxGetBugFieldOptions", params)
	dataStr, ok := client.Get(url)

	bugFields := model.ZentaoBugFields{}

	if ok {
		jsonData, err := simplejson.NewJson([]byte(dataStr))

		if err == nil {
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

		} else {
			logUtils.PrintToCmd(err.Error(), color.FgRed)
		}

		vari.ZenTaoBugFields = bugFields
	}
}

//func GetCaseModules(productId string) {
//	conf := configUtils.ReadCurrConfig()
//	Login(conf.Url, conf.Account, conf.Password)
//
//	params := [][]string{{"rootID", productId}, {"type", "case"}}
//	url := conf.Url + zentaoUtils.GenSuperApiUri("tree", "getOptionMenu", params)
//
//	dataStr, ok := client.Get(url, nil)
//
//	var moduelMap map[string]string
//	if ok {
//		json.Unmarshal([]byte(dataStr), &moduelMap)
//	}
//
//	vari.ZentaoCaseFileds.Modules = moduelMap
//}

func fieldMapToListOrderByInt(mp map[string]interface{}) []model.Option {
	arr := make([]model.Option, 0)

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

		opt := model.Option{Id: keyStr, Name: name}
		arr = append(arr, opt)
	}

	return arr
}

func fieldMapToListOrderByStr(mp map[string]interface{}, notNull bool) []model.Option {
	arr := make([]model.Option, 0)

	keys := make([]string, 0)
	for key, _ := range mp {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		name := strings.TrimSpace(mp[key].(string))
		if name == "" {
			if notNull {
				continue
			}
		}

		opt := model.Option{Id: key, Name: name}
		arr = append(arr, opt)
	}

	return arr
}

func fieldArrToListKeyStr(arr0 []interface{}, notNull bool) []model.Option {
	arr := make([]model.Option, 0)

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

		opt := model.Option{Id: val.(string), Name: name}
		arr = append(arr, opt)
	}

	return arr
}

func GetFirstNoEmptyVal(options []model.Option) string {
	for _, opt := range options {
		if opt.Name != "" {
			return opt.Id
		}
	}

	return ""
}

func GetIdByName(name string, options []model.Option) string {
	for _, opt := range options {
		if opt.Name == name {
			return opt.Id
		}
	}

	return ""
}

func GetNameById(id string, options []model.Option) string {
	for _, opt := range options {
		if opt.Id == id {
			return opt.Name
		}
	}

	return ""
}
