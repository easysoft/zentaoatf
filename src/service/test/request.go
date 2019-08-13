package testingService

import (
	httpClient "github.com/easysoft/zentaoatf/src/client"
	"github.com/easysoft/zentaoatf/src/utils/common"
	config2 "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	print2 "github.com/easysoft/zentaoatf/src/utils/print"
)

func GetZentaoSettings() {
	config := config2.ReadCurrConfig()

	entityType := config.EntityType
	entityVal := config.EntityVal

	requestObj := make(map[string]string)
	requestObj["entityType"] = entityType
	requestObj["entityVal"] = entityVal

	url := config.Url
	url = commonUtils.UpdateUrl(url)
	_, _ = httpClient.Post(url+constant.UrlZentaoSettings, requestObj)
	print2.PrintToCmd(url + constant.UrlZentaoSettings)

	//if err == nil {
	//	if pass {
	//		utils.PrintToCmd("success to get settings")
	//		//utils.ZendaoSettings = body.ZentaoSettings
	//	}
	//} else {
	//	utils.PrintToCmd(err.Error())
	//}
}
