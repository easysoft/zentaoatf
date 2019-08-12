package zentao

import (
	"github.com/easysoft/zentaoatf/src/http"
)

func ListCaseByProduct(baseUrl string, productId string) []byte {
	url := baseUrl + "case.json"

	params := make(map[string]string)
	params["productId"] = productId

	_, _ = http.Get(url, params)

	return nil
}

func ListCaseByTask(baseUrl string, taskId string) []byte {
	url := baseUrl + "case.json"

	params := make(map[string]string)
	params["taskId"] = taskId

	_, _ = http.Get(url, params)

	return nil
}
