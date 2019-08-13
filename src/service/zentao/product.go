package zentaoService

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
)

func GetProductInfo(baseUrl string, productId string) model.Product {
	params := [][]string{{"productID", productId}}

	myurl := baseUrl + zentaoUtils.GenSuperApiUri("product", "getById", params)
	dataStr, ok := client.Get(myurl, nil)

	if ok {
		var product model.Product
		json.Unmarshal([]byte(dataStr), &product)

		return product
	}

	return model.Product{}
}
