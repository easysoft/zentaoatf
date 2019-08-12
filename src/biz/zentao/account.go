package zentao

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/utils"
)

func Login(baseUrl string, account string, password string) {
	GetConfig(baseUrl)

	url := baseUrl + "user-login"

	params := make(map[string]string)
	params["account"] = account
	params["password"] = password

	_, err := http.Post(url, params)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetConfig(baseUrl string) {
	url := baseUrl + "?mode=getconfig"

	body, err := http.Get(url, nil)

	if err == nil {
		json, _ := simplejson.NewJson([]byte(body))

		utils.SessionId, _ = json.Get("sessionID").String()
		utils.SessionVar, _ = json.Get("sessionVar").String()

		fmt.Sprintf("%s: %s", utils.SessionVar, utils.SessionId)
	}
}

func GetSession(baseUrl string) {
	url := baseUrl + "api-getsessionid.json"

	body, _ := http.Get(url, nil)

	json, _ := simplejson.NewJson([]byte(body))
	status, _ := json.Get("status").String()

	pass := status == "" || status == "success" // some api not return a status

	if pass {
		dataStr, _ := json.Get("data").String()
		data, _ := simplejson.NewJson([]byte(dataStr))

		utils.SessionId, _ = data.Get("sessionID").String()
		utils.SessionVar, _ = data.Get("sessionName").String()

		fmt.Sprintf("%s: %s", utils.SessionVar, utils.SessionId)
	}
}
