package zentao

import (
	"github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/utils"
)

func GetTaskInfo(baseUrl string, productId string) []byte {
	url := baseUrl + "task.json"

	params := make(map[string]string)
	params["productId"] = productId

	json, _ := http.Get(url, params)

	return json
}

func GetCurrTaskInfo() []byte {
	conf := utils.ReadCurrConfig()

	baseUrl := conf.Url
	productId := conf.EntityVal

	json := GetTaskInfo(baseUrl, productId)

	return json
}
