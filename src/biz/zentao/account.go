package zentao

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/utils"
)

func Login() {
	conf := utils.ReadCurrConfig()
	url := conf.Url + "api-getsessionid.json"

	account := conf.Account
	password := conf.Password

	params := make(map[string]string)
	params["account"] = account
	params["password"] = password

	_, json, _ := http.Get(url, params)
	fmt.Println(json.String())
}

func GetSession() {
	conf := utils.ReadCurrConfig()
	url := conf.Url + "api-getsessionid.json"

	pass, json, _ := http.Get(url, nil)

	if pass {
		sessionID, _ := json.Get("sessionID").String()

		fmt.Println("sessionID: " + sessionID)
	}
}
