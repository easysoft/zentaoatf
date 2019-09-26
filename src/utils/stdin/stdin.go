package stdinUtils

import (
	"bufio"
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

	var numb string

	moduleCheck := ""
	productCheck := ""
	suiteCheck := ""
	taskCheck := ""

	if *moduleId != "" {
		moduleCheck = "*"
		numb = "2"
	} else if *productId != "" {
		productCheck = "*"
		numb = "1"
	} else if *suiteId != "" {
		suiteCheck = "*"
		numb = "3"
	} else if *taskId != "" {
		taskCheck = "*"
		numb = "4"
	}

	coType := GetInput("(1|2|3|4)", numb, "enter_co_type", productCheck, moduleCheck, suiteCheck, taskCheck)

	coType = strings.ToLower(coType)
	if coType == "1" {
		*productId = GetInput("\\d+", *productId,
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("product_id")+": "+*productId)

	} else if coType == "2" {
		*productId = GetInput("\\d+", *productId,
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("product_id")+": "+*productId)

		*moduleId = GetInput("\\d+", *moduleId,
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("module_id")+": "+*moduleId)

	} else if coType == "3" {
		*suiteId = GetInput("\\d+", *suiteId,
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("suite_id")+": "+*suiteId)
	} else if coType == "4" {
		*taskId = GetInput("\\d+", *taskId,
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("task_id")+": "+*taskId)
	}

	InputForBool(independentFile, false, "enter_co_independent")

	numbs, names, labels := langUtils.GetSupportLanguageOptions(nil)
	fmtParam := strings.Join(labels, "\n")

	langStr := GetInput("("+strings.Join(numbs, "|")+")", "", "enter_co_language", fmtParam)
	langNumb, _ := strconv.Atoi(langStr)

	*scriptLang = names[langNumb-1]
}

func InputForDir(dir *string, entity string) {
	*dir = GetInput("is_dir", "", "enter_dir", i118Utils.I118Prt.Sprintf(entity))
}

func InputForBool(in *bool, defaultVal bool, fmtStr string, fmtParam ...interface{}) {
	str := GetInput("(yes|no|y|n|)", "", fmtStr, fmtParam...)

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
		// fmt.Scanln(&ret)
		Scanf(&ret)

		if strings.TrimSpace(ret) == "" && defaultVal != "" {
			ret = defaultVal

			logUtils.PrintToStdOut(ret, -1)
		}

		temp := strings.ToLower(ret)
		if temp == "exit" {
			os.Exit(1)
		}

		logUtils.PrintToStdOut(ret, -1)

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

func Scanf(a *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	*a = string(data)
}
