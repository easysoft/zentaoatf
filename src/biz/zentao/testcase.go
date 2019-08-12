package zentao

import (
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/utils"
)

func ListCaseByProduct(baseUrl string, productId string) *simplejson.Json {
	params := map[string]string{"productID": productId}
	url := baseUrl + utils.GenSuperApiUri("testcase", "getModuleCases", params)
	body, err := http.Post(url, nil)

	if err == nil {
		json, _ := simplejson.NewJson([]byte(body))

		status, _ := json.Get("status").String()
		if status == "success" {
			dataStr, _ := json.Get("data").String()
			data, _ := simplejson.NewJson([]byte(dataStr))

			return data
		}
	}

	return nil
}

func ListCaseByTask(baseUrl string, taskId string) []byte {
	url := baseUrl + "case.json"

	params := make(map[string]string)
	params["taskId"] = taskId

	_, _ = http.Get(url, params)

	return nil
}
