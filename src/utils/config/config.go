package configUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/display"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	stdinUtils "github.com/easysoft/zentaoatf/src/utils/stdin"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
)

func InitConfig() {
	// preference from yaml
	vari.Config = getInst()

	// screen size
	InitScreenSize()

	// internationalization
	i118Utils.InitI118(vari.Config.Language)
}

func InitScreenSize() {
	w, h := display.GetScreenSize()
	vari.ScreenWidth = w
	vari.ScreenHeight = h
}

func SaveConfig(language string, url string, account string, password string) error {
	config := ReadCurrConfig()

	config.Version = constant.ConfigVer

	if language != "" {
		config.Language = language
	}
	if url != "" {
		config.Url = url
	}
	if account != "" {
		config.Account = account
	}
	if password != "" {
		config.Password = password
	}

	data, _ := yaml.Marshal(&config)
	ioutil.WriteFile(constant.ConfigFile, data, 0666)

	vari.Config = ReadCurrConfig()
	return nil
}

func PrintCurrConfig() {
	logUtils.PrintToStdOut("\n"+i118Utils.I118Prt.Sprintf("current_config"), color.FgCyan)

	val := reflect.ValueOf(vari.Config)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(vari.Config).NumField(); i++ {
		val := val.Field(i)
		name := typeOfS.Field(i).Name

		fmt.Printf("  %s: %v \n", name, val.Interface())
	}
}

func ReadCurrConfig() model.Config {
	config := model.Config{}

	configPath := constant.ConfigFile

	if !fileUtils.FileExist(configPath) {
		config.Language = "en"
		i118Utils.InitI118("en")

		return config
	}

	buf, _ := ioutil.ReadFile(configPath)
	yaml.Unmarshal(buf, &config)

	config.Url = commonUtils.UpdateUrl(config.Url)

	return config
}

func getInst() model.Config {
	CheckConfig()

	vari.Config = model.Config{}

	buf, _ := ioutil.ReadFile(constant.ConfigFile)
	yaml.Unmarshal(buf, &vari.Config)

	if vari.Config.Version != constant.ConfigVer { // old config file, re-init
		if vari.Config.Language != "en" && vari.Config.Language != "zh" {
			vari.Config.Language = "en"
		}

		SaveConfig(vari.Config.Language, vari.Config.Url, vari.Config.Account, vari.Config.Password)
	}

	return vari.Config
}

func CheckConfig() {
	configPath := constant.ConfigFile
	if !fileUtils.FileExist(configPath) {
		InputForSet()
	}
}

func InputForSet() {
	conf := ReadCurrConfig()

	var configSite bool

	language := conf.Language
	url := conf.Url
	account := conf.Account
	password := conf.Password

	logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("begin_config"), color.FgCyan)

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

	SaveConfig(language, url, account, password)

	PrintCurrConfig()
}

func CheckRequestConfig() {
	conf := ReadCurrConfig()
	if conf.Url == "" || conf.Account == "" || conf.Password == "" {
		InputForRequest()
	}
}

func InputForRequest() {
	conf := ReadCurrConfig()

	url := ""
	account := ""
	password := ""

	logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("need_config"), color.FgCyan)

	url = stdinUtils.GetInput("(http://.*)", conf.Url, "enter_url", conf.Url)

	account = stdinUtils.GetInput("(.{2,})", conf.Account, "enter_account", conf.Account)

	password = stdinUtils.GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)

	SaveConfig("", url, account, password)
}
