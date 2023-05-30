package zentaoHelper

import (
	"fmt"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	"path"
)

const (
	ApiPath = "api.php/v1/"
)

func GenApiUrl(pth string, params map[string]interface{}, baseUrl string) (url string) {
	uri := path.Join(ApiPath, pth)

	index := 0
	for key, val := range params {
		if index == 0 {
			uri += "?"
		} else {
			uri += "&"
		}

		uri += fmt.Sprintf("%v=%v", key, val)
		index++
	}

	url = httpUtils.AddSepIfNeeded(baseUrl) + uri

	return
}
