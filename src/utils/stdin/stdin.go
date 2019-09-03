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

func InputForSet() {
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

func CheckRequestConfig() {
	conf := configUtils.ReadCurrConfig()
	if conf.Url == "" || conf.Account == "" || conf.Password == "" {
		InputForRequest()
	}
}

func InputForRequest() {
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

func InputForCheckout(productId *string, moduleId *string, suiteId *string, taskId *string,
	independentFile *bool, scriptLang *string) {

	coType := getInput("(1|2|3|4)", "enter_co_type")

	coType = strings.ToLower(coType)
	if coType == "1" {
		*productId = getInput("\\d+",
			i118Utils.I118Prt.Sprintf("pls_input")+" "+i118Utils.I118Prt.Sprintf("product_id"))

	} else if coType == "2" {
		*productId = getInput("\\d+",
			i118Utils.I118Prt.Sprintf("pls_input")+" "+i118Utils.I118Prt.Sprintf("product_id"))

		*moduleId = getInput("\\d+",
			i118Utils.I118Prt.Sprintf("pls_input")+" "+i118Utils.I118Prt.Sprintf("module_id"))

	} else if coType == "3" {
		*suiteId = getInput("\\d+",
			i118Utils.I118Prt.Sprintf("pls_input")+" "+i118Utils.I118Prt.Sprintf("suite_id"))
	} else if coType == "4" {
		*taskId = getInput("\\d+",
			i118Utils.I118Prt.Sprintf("pls_input")+" "+i118Utils.I118Prt.Sprintf("task_id"))
	}

	InputForBool(independentFile, false, "enter_co_independent")

	regx := langUtils.GetSupportLangageRegx()
	fmtParam := strings.Join(langUtils.GetSupportLangageArr(), " / ")
	*scriptLang = getInput(regx, "enter_co_language", fmtParam)
}

func InputForDir(dir *string, entity string) {
	*dir = getInput("is_dir", "enter_dir", i118Utils.I118Prt.Sprintf(entity))
}

func InputForInt(in *string, fmtStr string, fmtParam ...string) {
	*in = getInput("\\d+", "fmtStr", fmtParam)
}

func InputForBool(in *bool, defaultVal bool, fmtStr string, fmtParam ...string) {
	str := getInput("(yes|no|y|n|)", fmtStr, fmtParam)

	if str == "" {
		*in = defaultVal
		return
	}

	if str == "y" && str != "yes" {
		*in = true
	} else {
		*in = false
	}
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
			ret = ""
			logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf(msg), color.FgRed)
		}
	}
}
