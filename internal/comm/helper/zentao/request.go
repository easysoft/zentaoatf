package zentaoHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"path"
	"strings"
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

	url = baseUrl + uri

	return
}

func GenApiUriOld(module string, methd string, param string) string {
	var uri string

	if commConsts.RequestType == commConsts.PathInfo {
		uri = fmt.Sprintf("%s-%s-%s.json", module, methd, param)
	} else {
		uri = fmt.Sprintf("index.php?m=%s&f=%s&%s&t=json", module, methd, param)
	}

	return uri
}

func GetRespErr(bytes []byte, key string) (err error) {
	if strings.Index(string(bytes), "login") > -1 {
		err = errors.New(i118Utils.Sprintf("fail_to_login"))
		return
	}

	var respData = serverDomain.ZentaoRespData{}
	err = json.Unmarshal(bytes, &respData)

	// map[result:success] or map[status:success, data:{}]
	if err == nil && (respData.Result != "" || respData.Result != "success") {
		msg := i118Utils.Sprintf(key, respData.Message)
		err = errors.New(msg)
	}

	return
}
