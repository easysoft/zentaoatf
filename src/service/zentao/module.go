package zentaoService

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
)

func ListCaseModule(baseUrl string, productId string) []model.Module {
	params := [][]string{{"rootID", productId}, {"type", "case"}}

	myurl := baseUrl + zentaoUtils.GenSuperApiUri("tree", "getOptionMenu", params)
	dataStr, ok := client.Get(myurl, nil)

	modules := make([]model.Module, 0)
	if ok {
		var moduleMap map[int]string
		json.Unmarshal([]byte(dataStr), &moduleMap)

		for id, name := range moduleMap {
			modules = append(modules, model.Module{Id: id, Name: name})
		}

		return modules
	}

	return nil
}
