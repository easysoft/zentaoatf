package stdinUtils

import (
	"fmt"
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

func InputForCheckout(productId *string, moduleId *string, suiteId *string, taskId *string,
	independentFile *bool, scriptLang *string) {

	coType := GetInput("(1|2|3|4)", "", "enter_co_type")

	coType = strings.ToLower(coType)
	if coType == "1" {
		*productId = GetInput("\\d+", "",
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("product_id"))

	} else if coType == "2" {
		*productId = GetInput("\\d+", "",
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("product_id"))

		*moduleId = GetInput("\\d+", "",
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("module_id"))

	} else if coType == "3" {
		*suiteId = GetInput("\\d+", "",
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("suite_id"))
	} else if coType == "4" {
		*taskId = GetInput("\\d+", "",
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("task_id"))
	}

	InputForBool(independentFile, false, "enter_co_independent")

	numbs, names, labels := langUtils.GetSupportLanguageOptions()
	fmtParam := strings.Join(labels, "\n")
	numbStr := GetInput("("+strings.Join(numbs, "|")+")", "enter_co_language", fmtParam)

	numb, _ := strconv.Atoi(numbStr)

	*scriptLang = names[numb-1]
}

func InputForDir(dir *string, entity string) {
	*dir = GetInput("is_dir", "", "enter_dir", i118Utils.I118Prt.Sprintf(entity))
}

func InputForBool(in *bool, defaultVal bool, fmtStr string, fmtParam ...string) {
	str := GetInput("(yes|no|y|n|)", "", fmtStr, fmtParam)

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

func GetInput(regx string, defaultVal string, fmtStr string, params ...interface{}) string {
	var ret string

	msg := i118Utils.I118Prt.Sprintf(fmtStr, params...)

	for {
		logUtils.PrintToStdOut("\n"+msg, color.FgCyan)
		fmt.Scanln(&ret)

		if strings.TrimSpace(ret) == "" {
			ret = defaultVal
		}

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
