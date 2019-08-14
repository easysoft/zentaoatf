package zentaoService

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/service/client"
	"github.com/easysoft/zentaoatf/src/utils/vari"
)

func Login(baseUrl string, account string, password string) {
	GetConfig(baseUrl)

	url := baseUrl + "user-login"

	params := make(map[string]string)
	params["account"] = account
	params["password"] = password

	_, ok := client.PostStr(url, params)
	if ok {
		fmt.Println("succes to login")
	} else {
		fmt.Println("fail to login")
	}
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

		fmt.Sprintf("%s: %s", vari.SessionVar, vari.SessionId)
	}
}

func GetSession(baseUrl string) {
	url := baseUrl + "api-getsessionid.json"

	dataStr, ok := client.Get(url, nil)
	if ok {
		data, _ := simplejson.NewJson([]byte(dataStr))

		vari.SessionId, _ = data.Get("sessionID").String()
		vari.SessionVar, _ = data.Get("sessionName").String()

		fmt.Sprintf("%s: %s", vari.SessionVar, vari.SessionId)
	}
}
