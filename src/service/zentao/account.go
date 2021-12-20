package zentaoService

import (
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/service/client"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"strings"
)

func Login(baseUrl string, account string, password string) bool {
	ok := GetConfig(baseUrl)

	if !ok {
		logUtils.PrintToCmd(i118Utils.Sprintf("fail_to_login"), color.FgRed)
		return false
	}

	uri := ""
	if vari.RequestType == constant.RequestTypePathInfo {
		uri = "user-login.json"
	} else {
		uri = "index.php?m=user&f=login&t=json"
	}
	url := baseUrl + uri

	params := make(map[string]string)
	params["account"] = account
	params["password"] = password

	var body string
	body, ok = client.PostStr(url, params)
	if !ok || (ok && strings.Index(body, "title") > 0) { // use PostObject to login again for new system
		_, ok = client.PostObject(url, params, true)
	}
	if ok {
		if vari.Verbose {
			logUtils.Screen(i118Utils.Sprintf("success_to_login"))
		}
	} else {
		logUtils.PrintToCmd(i118Utils.Sprintf("fail_to_login"), color.FgRed)
	}

	return ok
}

func GetConfig(baseUrl string) bool {
	if vari.RequestType != "" {
		return true
	}

	// get config
	url := baseUrl + "?mode=getconfig"
	body, ok := client.Get(url)
	if !ok {
		return false
	}

	json, _ := simplejson.NewJson([]byte(body))
	vari.ZenTaoVersion, _ = json.Get("version").String()
	vari.SessionId, _ = json.Get("sessionID").String()
	vari.SessionVar, _ = json.Get("sessionVar").String()
	vari.RequestType, _ = json.Get("requestType").String()
	vari.RequestFix, _ = json.Get("requestFix").String()

	// check site path by calling login interface
	uri := ""
	if vari.RequestType == constant.RequestTypePathInfo {
		uri = "user-login.json"
	} else {
		uri = "index.php?m=user&f=login&t=json"
	}
	url = baseUrl + uri
	body, ok = client.Get(url)
	if !ok {
		return false
	}

	return true
}
