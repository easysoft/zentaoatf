package biz

import (
	httpClient "github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/utils"
)

func GetZentaoSettings() {
	config := utils.ReadCurrConfig()

	entityType := config.EntityType
	entityVal := config.EntityVal

	requestObj := make(map[string]string)
	requestObj["entityType"] = entityType
	requestObj["entityVal"] = entityVal

	url := config.Url
	url = utils.UpdateUrl(url)
	_, _ = httpClient.Post(url+utils.UrlZentaoSettings, requestObj)
	utils.PrintToCmd(url + utils.UrlZentaoSettings)

	//if err == nil {
	//	if pass {
	//		utils.PrintToCmd("success to get settings")
	//		//utils.ZendaoSettings = body.ZentaoSettings
	//	}
	//} else {
	//	utils.PrintToCmd(err.Error())
	//}
}
