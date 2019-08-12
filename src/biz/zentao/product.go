package zentao

import (
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/utils"
)

func GetProductInfo(baseUrl string, productId string) *simplejson.Json {
	url := baseUrl + "product-browse-" + productId + ".json"
	body, err := http.Get(url, nil)

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

func GetCurrProductInfo() *simplejson.Json {
	conf := utils.ReadCurrConfig()

	baseUrl := conf.Url
	productId := conf.EntityVal

	json := GetProductInfo(baseUrl, productId)

	return json
}
