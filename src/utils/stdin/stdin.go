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
	"strconv"
	"strings"
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

	language = getInput("(1|2|)", "enter_language", enCheck, zhCheck)

	if language == "" {
		language = conf.Language
		logUtils.PrintToStdOut(numb, -1)
	} else if language == "1" {
		language = "en"
	} else {
		language = "zh"
	}

	InputForBool(&configSite, true, "config_zentao_site")
	if configSite {
		url = getInput("(http://.*|)", "enter_url", conf.Url)
		if url == "" {
			url = conf.Url
			logUtils.PrintToStdOut(url, -1)
		}

		account = getInput("(.{3,}|)", "enter_account", conf.Account)
		if account == "" {
			account = conf.Account
			logUtils.PrintToStdOut(account, -1)
		}

		password = getInput("(.{4,}|)", "enter_password", conf.Password)
		if password == "" {
			password = conf.Password
			logUtils.PrintToStdOut(password, -1)
		}

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
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("product_id"))

	} else if coType == "2" {
		*productId = getInput("\\d+",
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("product_id"))

		*moduleId = getInput("\\d+",
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("module_id"))

	} else if coType == "3" {
		*suiteId = getInput("\\d+",
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("suite_id"))
	} else if coType == "4" {
		*taskId = getInput("\\d+",
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("task_id"))
	}

	InputForBool(independentFile, false, "enter_co_independent")

	numbs, names, labels := langUtils.GetSupportLanguageOptions()
	fmtParam := strings.Join(labels, "\n")
	numbStr := getInput("("+strings.Join(numbs, "|")+")", "enter_co_language", fmtParam)

	numb, _ := strconv.Atoi(numbStr)

	*scriptLang = names[numb-1]
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

		msg := ""
		if *in {
			msg = "Yes"
		} else {
			msg = "No"
		}

		logUtils.PrintToStdOut(msg, -1)
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
