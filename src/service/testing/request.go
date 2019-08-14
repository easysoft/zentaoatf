package testingService

import (
	config2 "github.com/easysoft/zentaoatf/src/utils/config"
)

func GetZentaoSettings() {
	config := config2.ReadCurrConfig()

	entityType := config.EntityType
	entityVal := config.EntityVal

	requestObj := make(map[string]interface{})
	requestObj["entityType"] = entityType
	requestObj["entityVal"] = entityVal

	//url := config.Url
	//url = commonUtils.UpdateUrl(url)
	//_, _ = httpClient.PostJson(url+constant.UrlZentaoSettings, requestObj)
	//print2.PrintToCmd(url + constant.UrlZentaoSettings)

	//if err == nil {
	//	if pass {
	//		utils.PrintToCmd("success to get settings")
	//		//utils.ZendaoSettings = body.ZentaoSettings
	//	}
	//} else {
	//	utils.PrintToCmd(err.Error())
	//}
}
