package zentao

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/src/client"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
)

func GetProductInfo(baseUrl string, productId string) model.Product {
	params := map[string]string{"productID": productId}

	myurl := baseUrl + utils.GenSuperApiUri("product", "getById", params)
	dataStr, ok := client.Get(myurl, nil)

	if ok {
		var product model.Product
		json.Unmarshal([]byte(dataStr), &product)

		return product
	}

	return model.Product{}
}
