package httpHelper

import (
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	"github.com/tidwall/gjson"
)

func Login() (ret string) {
	url := zentaoHelper.GenApiUrl("tokens", nil, constTestHelper.ZtfUrl)

	params := map[string]string{
		"account":  constTestHelper.ZentaoUsername,
		"password": constTestHelper.ZentaoPassword,
	}
	bodyBytes, err := Post(url, "", params)
	if err != nil {
		return
	}

	ret = gjson.Get(string(bodyBytes), "token").String()

	return
}
