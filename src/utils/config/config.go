package configUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/display"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/lang"
	"github.com/easysoft/zentaoatf/src/utils/log"
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

func ConfigForSet() {
	configSite := ""

	language := ""
	url := ""
	account := ""
	password := ""

	fmt.Println(i118Utils.I118Prt.Sprintf("begin_config"))

	language = getInput("(english|chinese|e|c|)", "enter_language")
	languageDefault := "en"
	if language == "chinese" || language == "c" { // default en
		language = "zh"
	} else {
		if language == "" {
			fmt.Print("English")
		}
		language = languageDefault
	}

	configSite = getInput("(yes|no|y|n|)", "config_zentao_site")
	configSiteDefault := "yes"
	if configSite != "no" && configSite != "n" { // default yes
		if configSite == "" {
			fmt.Print(configSiteDefault)
		}

		url = getInput("http://.*", "enter_url")
		account = getInput(".{3,}", "enter_account")
		password = getInput(".{4,}", "enter_password")
	}

	SaveConfig(language, url, account, password)

	PrintCurrConfig()
}

func ConfigForCheckout(productId *string, moduleId *string, suiteId *string, taskId *string,
	independentFile *bool, scriptLang *string) {

	coType := getInput("(product|module|suite|task|p|m|s|t)", "enter_co_type")

	coType = strings.ToLower(coType)
	if coType == "product" || coType == "p" {
		*productId = getInput("\\d+", "productId")
	} else if coType == "module" || coType == "m" {
		*productId = getInput("\\d+", "productId")
		*moduleId = getInput("\\d+", "moduleId")
	} else if coType == "suite" || coType == "s" {
		*suiteId = getInput("\\d+", "suiteId")
	} else if coType == "task" || coType == "t" {
		*taskId = getInput("\\d+", "taskId")
	}

	indep := getInput("(yes|no|y|n|)", "enter_co_independent")
	indep = strings.ToLower(indep)
	if indep == "yes" && indep == "y" {
		*independentFile = true
	} else {
		if indep == "" {
			fmt.Print("no")
		}
		*independentFile = false
	}

	regx := "(" + strings.Join(langUtils.GetSupportLangageArr(), "|") + ")"
	fmtParam := strings.Join(langUtils.GetSupportLangageArr(), " / ")
	*scriptLang = getInput(regx, "enter_co_language", fmtParam)

	PrintCurrConfig()
}

func ConfigForDir(dir *string, entity string) {
	*dir = getInput("is_dir", "enter_dir", i118Utils.I118Prt.Sprintf(entity))
}

func ConfigForInt(in *string, entity string) {
	*in = getInput("\\d+", "enter_id", i118Utils.I118Prt.Sprintf(entity))
}

func getInput(regx string, fmtStr string, params ...interface{}) string {
	var ret string

	msg := i118Utils.I118Prt.Sprintf(fmtStr, params...)

	for {
		color.Cyan("\n" + msg)
		fmt.Scanln(&ret)

		temp := strings.ToLower(ret)
		if temp == "exit" {
			os.Exit(1)
		}

		if regx == "" {
			return ret
		}

		var pass bool
		var msg string
		if regx == "is_dir" {
			pass = fileUtils.IsDir(ret)
			msg = "dir_not_exist"
		} else {
			pass, _ = regexp.MatchString("^"+regx+"$", temp)
			msg = "invalid_input"
		}

		if pass {
			return ret
		} else {
			color.Red(i118Utils.I118Prt.Sprintf(msg) + "\n")
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
	w, h := display.GetScreenSize()
	vari.ScreenWidth = w
	vari.ScreenHeight = h
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
			saveEmptyConfig()
		}
	})
	return vari.Config
}

func SaveConfig(language string, url string, account string, password string) error {
	config := ReadCurrConfig()

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

func ReadCurrConfig() model.Config {
	configPath := constant.ConfigFile
	var config model.Config

	if !fileUtils.FileExist(configPath) {
		saveEmptyConfig()
	}
	buf, _ := ioutil.ReadFile(configPath)
	yaml.Unmarshal(buf, &config)

	config.Url = commonUtils.UpdateUrl(config.Url)

	return config
}

func saveEmptyConfig() error {
	config := model.Config{Language: "en", Url: "", Account: "", Password: ""}

	data, _ := yaml.Marshal(&config)
	ioutil.WriteFile(constant.ConfigFile, data, 0666)

	return nil
}
