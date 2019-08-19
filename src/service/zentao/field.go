package zentaoService

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"sort"
	"strconv"
)

func GetBugFiledOptions() {
	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	params := fmt.Sprintf("%d-%d", conf.ProductId, conf.ProjectId)
	url := conf.Url + zentaoUtils.GenApiUri("bug", "ajaxGetBugFieldOptions", params)
	dataStr, ok := client.Get(url, nil)

	bugFields := model.ZentaoBugFileds{}

	if ok {
		jsonData, err := simplejson.NewJson([]byte(dataStr))

		if err == nil {
			mp, _ := jsonData.Get("modules").Map()
			bugFields.Modules = fieldMapToListKeyInt(mp)

			mp, _ = jsonData.Get("categories").Map()
			bugFields.Categories = fieldMapToListKeyStr(mp)

			mp, _ = jsonData.Get("versions").Map()
			bugFields.Versions = fieldMapToListKeyStr(mp)

			mp, _ = jsonData.Get("severities").Map()
			bugFields.Severities = fieldMapToListKeyInt(mp)

			arr, _ := jsonData.Get("priorities").Array()
			bugFields.Priorities = fieldArrToListKeyStr(arr)

		} else {
			logUtils.PrintToCmd(err.Error())
		}

		vari.ZentaoBugFileds = bugFields
	}
}

func fieldMapToListKeyInt(mp map[string]interface{}) []model.Option {
	arr := make([]model.Option, 0)

	keys := make([]int, 0)
	for key, _ := range mp {
		keyint, _ := strconv.Atoi(key)
		keys = append(keys, keyint)
	}

	sort.Ints(keys)

	for _, key := range keys {
		keyStr := strconv.Itoa(key)

		opt := model.Option{Id: keyStr, Name: mp[keyStr].(string)}
		arr = append(arr, opt)
	}

	return arr
}

func fieldMapToListKeyStr(mp map[string]interface{}) []model.Option {
	arr := make([]model.Option, 0)

	keys := make([]string, 0)
	for key, _ := range mp {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		opt := model.Option{Id: key, Name: mp[key].(string)}
		arr = append(arr, opt)
	}

	return arr
}

func fieldArrToListKeyStr(arr0 []interface{}) []model.Option {
	arr := make([]model.Option, 0)

	keys := make([]string, 0)
	for _, val := range arr0 {
		keys = append(keys, val.(string))
	}

	sort.Strings(keys)

	for _, val := range arr0 {
		opt := model.Option{Id: val.(string), Name: val.(string)}
		arr = append(arr, opt)
	}

	return arr
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
