package configUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	assertUtils "github.com/easysoft/zentaoatf/src/utils/assert"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/display"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	stdinUtils "github.com/easysoft/zentaoatf/src/utils/stdin"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"gopkg.in/ini.v1"
	"os"
	"reflect"
)

func InitConfig() {
	vari.ZtfDir = fileUtils.GetZtfDir()
	constant.ConfigFile = vari.ZtfDir + constant.ConfigFile

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

func SaveConfig(conf model.Config) error {
	fmt.Printf("\n%s=%s\n", "php", conf.Php)

	fileUtils.MkDirIfNeeded(fileUtils.GetZtfDir() + "conf")

	conf.Version = constant.ConfigVer

	cfg := ini.Empty()
	cfg.ReflectFrom(&conf)

	cfg.SaveTo(constant.ConfigFile)

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

	ini.MapTo(&config, constant.ConfigFile)

	config.Url = commonUtils.UpdateUrl(config.Url)

	return config
}

func getInst() model.Config {
	CheckConfig()

	configFile := constant.ConfigFile

	ini.MapTo(&vari.Config, configFile)

	if vari.Config.Version != constant.ConfigVer { // old config file, re-init
		if vari.Config.Language != "en" && vari.Config.Language != "zh" {
			vari.Config.Language = "en"
		}

		SaveConfig(vari.Config)
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

	//language := conf.Language
	//url := conf.Url
	//account := conf.Account
	//password := conf.Password

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
		conf.Language = "en"
	} else {
		conf.Language = "zh"
	}

	stdinUtils.InputForBool(&configSite, true, "config_zentao_site")
	if configSite {
		conf.Url = stdinUtils.GetInput("(http://.*)", conf.Url, "enter_url", conf.Url)

		conf.Account = stdinUtils.GetInput("(.{2,})", conf.Account, "enter_account", conf.Account)

		conf.Password = stdinUtils.GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)
	}

	scripts := assertUtils.GetCaseByDirAndFile([]string{"."})
	InputForScriptInterpreter(scripts, &conf)

	fmt.Printf("\n%s=%s\n", "php", conf.Php)

	SaveConfig(conf)
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

	logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("need_config"), color.FgCyan)

	conf.Url = stdinUtils.GetInput("(http://.*)", conf.Url, "enter_url", conf.Url)
	conf.Account = stdinUtils.GetInput("(.{2,})", conf.Account, "enter_account", conf.Account)
	conf.Password = stdinUtils.GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)

	SaveConfig(conf)
}

func InputForScriptInterpreter(scripts []string, config *model.Config) {
	//var configScriptInterpreter bool
	//stdinUtils.InputForBool(&configScriptInterpreter, true, "set_script_interpreter")
	//if configScriptInterpreter {
	langs := assertUtils.GetScriptType(scripts)

	for _, lang := range langs {
		inter := stdinUtils.GetInput(string(os.PathSeparator)+"+", commonUtils.GetFieldVal(*config, lang), "set_script_interpreter", lang)

		fmt.Printf("lang: %s, inter: %s", lang, inter)

		commonUtils.SetFieldVal(config, lang, inter)
	}

	fmt.Printf("\n%s=%s\n", "php", (*config).Php)

	//}
}
