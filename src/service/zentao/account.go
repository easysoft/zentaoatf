package zentaoService

import (
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/service/client"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
)

func Login(baseUrl string, account string, password string) bool {
	GetConfig(baseUrl)

	url := baseUrl + "user-login"

	params := make(map[string]string)
	params["account"] = account
	params["password"] = password

	var log string
	_, ok := client.PostStr(url, params)
	if ok {
		log = i118Utils.I118Prt.Sprintf("success_to_login")
		logUtils.PrintToCmd(log, -1)
	} else {
		log = i118Utils.I118Prt.Sprintf("fail_to_login")
		logUtils.PrintToCmd(log, color.FgRed)
	}

	return ok
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
	}
}
