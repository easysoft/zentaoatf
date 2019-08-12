package zentao

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
)

func GetProductInfo(baseUrl string, productId string) model.Product {
	params := map[string]string{"productID": productId}

	myurl := baseUrl + utils.GenSuperApiUri("product", "getById", params)
	bodyStr, err := http.Get(myurl, nil)

	if err == nil {
		var bodyJson model.ZentaoResponse
		json.Unmarshal(bodyStr, &bodyJson)

		if bodyJson.Status == "success" {
			dataStr := bodyJson.Data

			var product model.Product
			json.Unmarshal([]byte(dataStr), &product)

			return product
		}
	}

	return model.Product{}
}
