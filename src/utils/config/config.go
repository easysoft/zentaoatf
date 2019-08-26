package configUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/display"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

func InitConfig() {
	// preference from yaml
	vari.Config = getInst()

	// screen size
	InitScreenSize()

	// internationalization
	i118Utils.InitI118(vari.Config.Language)

	if strings.Index(os.Args[0], "atf") > -1 && (len(os.Args) > 1 && os.Args[1] != "set") {
		PrintCurrConfig()
	}
}

func ConfigFromStdin() {
	configSite := ""

	language := ""
	url := ""
	account := ""
	password := ""

	fmt.Printf(i118Utils.I118Prt.Sprintf("begin_config"))

	language = getInput("enter_language", "(english|chinese|e|c")
	if strings.Index(strings.ToLower(language), "e") == 0 {
		language = "en"
	} else {
		language = "zh"
	}

	configSite = getInput("config_zentao_site", "yes|no|y|n")
	if strings.Index(configSite, "y") != 0 {
		os.Exit(1)
	}

	url = getInput("enter_url", "http://.*")

	account = getInput("enter_account", ".[2,]")

	password = getInput("enter_password", ".[6,]")

}

func getInput(msg string, regx string) string {
	var ret string

	for {
		fmt.Printf(i118Utils.I118Prt.Sprintf(msg))
		fmt.Scanf("%s", &ret)

		ret = strings.ToLower(ret)
		if ret == "exit" {
			return ""
		}

		if regx == "" {
			return ret
		}

		pass, _ := regexp.MatchString(regx, ret)
		if pass {
			return ret
		}
	}
}

func SetLanguage(lang string, dumb bool) {
	buf, _ := ioutil.ReadFile(constant.ConfigFile)
	yaml.Unmarshal(buf, &vari.Config)

	vari.Config.Language = lang

	data, _ := yaml.Marshal(&vari.Config)
	ioutil.WriteFile(constant.ConfigFile, data, 0666)

	// re-init language resource
	i118Utils.InitI118(vari.Config.Language)

	if !dumb {
		logUtils.PrintToCmd(color.CyanString(i118Utils.I118Prt.Sprintf("set_config", i118Utils.I118Prt.Sprintf("language"),
			i118Utils.I118Prt.Sprintf(vari.Config.Language))))
	}
}

func PrintCurrConfig() {
	color.Cyan(i118Utils.I118Prt.Sprintf("current_config", ""))

	val := reflect.ValueOf(vari.Config)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(vari.Config).NumField(); i++ {
		val := val.Field(i)
		name := typeOfS.Field(i).Name

		if !vari.RunFromCui && (name == "Width" || name == "Height" || name == "WorkHistories") {
			continue
		}
		fmt.Printf("  %s: %v \n", typeOfS.Field(i).Name, val.Interface())
	}
}

func PrintConfigToView() {
	cmdView, _ := vari.Cui.View("cmd")
	fmt.Fprintln(cmdView, color.CyanString(i118Utils.I118Prt.Sprintf("current_config", "")))

	val := reflect.ValueOf(vari.Config)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(vari.Config).NumField(); i++ {
		val := val.Field(i)
		fmt.Fprintln(cmdView, fmt.Sprintf("  %s: %v", typeOfS.Field(i).Name, val.Interface()))
	}
}

func InitScreenSize() {
	w, _ := display.GetScreenSize()
	vari.ScreenWidth = w
}

func getInst() model.Config {
	var once sync.Once
	once.Do(func() {
		vari.Config = model.Config{}
		if fileUtils.FileExist(constant.ConfigFile) {
			buf, _ := ioutil.ReadFile(constant.ConfigFile)
			yaml.Unmarshal(buf, &vari.Config)
		} else { // init
			vari.Config.Language = "en"

			data, _ := yaml.Marshal(&vari.Config)
			ioutil.WriteFile(constant.ConfigFile, data, 0666)
		}
	})
	return vari.Config
}
