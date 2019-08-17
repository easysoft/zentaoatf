package zentaoService

import (
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/service/client"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
)

func Login(baseUrl string, account string, password string) {
	GetConfig(baseUrl)

	url := baseUrl + "user-login"

	params := make(map[string]string)
	params["account"] = account
	params["password"] = password

	var log string
	_, ok := client.PostStr(url, params)
	if ok {
		log = "succes to login"
	} else {
		log = "fail to login"
	}
	logUtils.PrintToCmd(log)
}

func GetConfig(baseUrl string) {
	url := baseUrl + "?mode=getconfig"

	body, ok := client.Get(url, nil)

	if ok {
		json, _ := simplejson.NewJson([]byte(body))

		vari.SessionId, _ = json.Get("sessionID").String()
		vari.SessionVar, _ = json.Get("sessionVar").String()
		vari.RequestType, _ = json.Get("requestType").String()
		vari.RequestFix, _ = json.Get("requestFix").String()

		//fmt.Sprintf("%s: %s", vari.SessionVar, vari.SessionId)
	}
}
