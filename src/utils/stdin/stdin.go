package stdinUtils

import (
	"fmt"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	langUtils "github.com/easysoft/zentaoatf/src/utils/lang"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/fatih/color"
	"os"
	"regexp"
	"strings"
)

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

	configUtils.SaveConfig(language, url, account, password)

	configUtils.PrintCurrConfig()
}

func CheckConfigForRequest() {
	conf := configUtils.ReadCurrConfig()
	if conf.Url == "" || conf.Account == "" || conf.Password == "" {
		ConfigForRequest()
	}
}

func ConfigForRequest() {
	url := ""
	account := ""
	password := ""

	fmt.Println(i118Utils.I118Prt.Sprintf("begin_config"))

	url = getInput("http://.*", "enter_url")
	account = getInput(".{3,}", "enter_account")
	password = getInput(".{4,}", "enter_password")

	configUtils.SaveConfig("", url, account, password)

	configUtils.PrintCurrConfig()
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

	configUtils.PrintCurrConfig()
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
		logUtils.PrintToStdOut("\n"+msg, color.FgCyan)
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
			logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf(msg), color.FgRed)
		}
	}
}
