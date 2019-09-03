package action

import (
	"fmt"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	stdinUtils "github.com/easysoft/zentaoatf/src/utils/stdin"
)

func InputForSet() {
	conf := configUtils.ReadCurrConfig()

	var configSite bool

	language := ""
	url := ""
	account := ""
	password := ""

	fmt.Println(i118Utils.I118Prt.Sprintf("begin_config"))

	enCheck := ""
	var numb string
	if conf.Language == "en" {
		enCheck = "*"
		numb = "1"
	}
	zhCheck := ""
	if conf.Language == "zh" {
		zhCheck = "*"
		numb = "2"
	}

	numbSelected := stdinUtils.GetInput("(1|2)", numb, "enter_language", enCheck, zhCheck)

	if numbSelected == "1" {
		language = "en"
	} else {
		language = "zh"
	}

	stdinUtils.InputForBool(&configSite, true, "config_zentao_site")
	if configSite {
		url = stdinUtils.GetInput("(http://.*)", conf.Url, "enter_url", conf.Url)

		account = stdinUtils.GetInput("(.{2,})", conf.Account, "enter_account", conf.Account)

		password = stdinUtils.GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)
	}

	configUtils.SaveConfig(language, url, account, password)

	configUtils.PrintCurrConfig()
}

func CheckRequestConfig() {
	conf := configUtils.ReadCurrConfig()
	if conf.Url == "" || conf.Account == "" || conf.Password == "" {
		InputForRequest()
	}
}

func InputForRequest() {
	conf := configUtils.ReadCurrConfig()

	url := ""
	account := ""
	password := ""

	fmt.Println(i118Utils.I118Prt.Sprintf("begin_config"))

	url = stdinUtils.GetInput("(http://.*)", conf.Url, "enter_url", conf.Url)

	account = stdinUtils.GetInput("(.{2,})", conf.Account, "enter_account", conf.Account)

	password = stdinUtils.GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)

	configUtils.SaveConfig("", url, account, password)
}
