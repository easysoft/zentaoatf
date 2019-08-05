package biz

import (
	"encoding/json"
	httpClient "github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/utils"
)

func GetZentaoSettings() {
	config := utils.ReadCurrConfig()

	entityType := config.EntityType
	entityVal := config.EntityVal

	requestObj := make(map[string]interface{})
	requestObj["entityType"] = entityType
	requestObj["entityVal"] = entityVal

	reqStr, _ := json.Marshal(requestObj)
	utils.PrintToCmd(string(reqStr))

	url := config.Url
	url = utils.UpdateUrl(url)
	body, err := httpClient.Post(url+utils.UrlZentaoSettings, string(reqStr))
	utils.PrintToCmd(url + utils.UrlZentaoSettings)

	if err == nil {
		if body.Code == 1 {
			utils.PrintToCmd("success to get settings")
			utils.ZendaoSettings = body.ZentaoSettings
		}
	} else {
		utils.PrintToCmd(err.Error())
	}
}
