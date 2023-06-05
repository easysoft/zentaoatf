package httpHelper

import (
	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	"github.com/tidwall/gjson"
)

func Login() (ret string) {
	url := zentaoHelper.GenApiUrl("tokens", nil, constTestHelper.ZentaoSiteUrl)

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
