package zentaoService

import (
	"github.com/easysoft/zentaoatf/src/model"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
)

func GetZentaoSettings() {
	config := configUtils.ReadCurrConfig()

	entityType := config.EntityType
	entityVal := config.EntityVal

	requestObj := make(map[string]interface{})
	requestObj["entityType"] = entityType
	requestObj["entityVal"] = entityVal

	//url := config.Url
	//_, _ = client.PostJson(url+constant.UrlZentaoSettings, requestObj)
	//printUtils.PrintToCmd(url + constant.UrlZentaoSettings)
	//
	//if err == nil {
	//	if pass {
	//		utils.PrintToCmd("success to get settings")
	//		//utils.ZendaoSettings = body.ZentaoSettings
	//	}
	//} else {
	//	printUtils.PrintToCmd(err.Error())
	//}
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
