package stdinUtils

import (
	"bufio"
	"fmt"
	langHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/lang"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"strconv"

	"github.com/fatih/color"
	"os"
	"regexp"
	"strings"
)

func InputForBool(in *bool, defaultVal bool, fmtStr string, fmtParam ...interface{}) {
	str := GetInput("(yes|no|y|n|)", "", fmtStr, fmtParam...)

	if str == "" {
		*in = defaultVal

		msg := ""
		if *in {
			msg = "yes"
		} else {
			msg = "no"
		}

		fmt.Print("\033[A")
		logUtils.ExecConsole(-1, msg)
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

	msg := i118Utils.Sprintf(fmtStr, params...)

	for {
		logUtils.ExecConsole(color.FgCyan, "\n"+msg)
		Scanf(&ret)
		ret = strings.TrimSpace(ret)

		if ret == "" && defaultVal != "" {
			ret = defaultVal
			logUtils.ExecConsole(-1, ret)
		}

		temp := strings.ToLower(ret)
		if temp == "exit" {
			color.Unset()
			os.Exit(0)
		}

		if regx == "" {
			return ret
		}

		var pass bool
		var msg string
		if regx == "is_dir" {
			pass = fileUtils.IsDir(ret)
			msg = "dir_not_exist"
		} else if regx != "" {
			pass, _ = regexp.MatchString("^"+regx+"$", temp)
			if !pass {
				msg = "invalid_input"
			}
		}

		if pass {
			return ret
		} else {
			ret = ""
			logUtils.ExecConsole(color.FgRed, i118Utils.Sprintf(msg))
		}
	}
}

func Scanf(a *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	*a = string(data)
}

func InputForCheckout(productId, moduleId, suiteId, taskId *string,
	independentFile *bool, scriptLang *string) {

	var numb string

	productCheckbox := ""
	suiteCheckbox := ""
	taskCheckbox := ""

	if *productId != "" {
		productCheckbox = "*"
		numb = "1"
	} else if *suiteId != "" {
		suiteCheckbox = "*"
		numb = "2"
	} else if *taskId != "" {
		taskCheckbox = "*"
		numb = "3"
	}

	coType := GetInput("(1|2|3)", numb, "enter_co_type", productCheckbox, suiteCheckbox, taskCheckbox)

	coType = strings.ToLower(coType)
	if coType == "1" {
		*productId = GetInput("\\d+", *productId,
			i118Utils.Sprintf("pls_enter")+" "+i118Utils.Sprintf("product_id")+": "+*productId)

		*moduleId = GetInput("\\d*", *moduleId,
			i118Utils.Sprintf("pls_enter")+" "+i118Utils.Sprintf("module_id")+": "+*moduleId)

	} else if coType == "2" {
		*suiteId = GetInput("\\d+", *suiteId,
			i118Utils.Sprintf("pls_enter")+" "+i118Utils.Sprintf("suite_id")+": "+*suiteId)
	} else if coType == "3" {
		*taskId = GetInput("\\d+", *taskId,
			i118Utils.Sprintf("pls_enter")+" "+i118Utils.Sprintf("task_id")+": "+*taskId)
	}

	InputForBool(independentFile, false, "enter_co_independent")

	numbs, names, labels := langHelper.GetSupportLanguageOptions(nil)
	fmtParam := make([]string, 0)
	dft := ""
	for idx, label := range labels {
		if names[idx] == *scriptLang {
			dft = strconv.Itoa(idx + 1)
			label += " *"
		}
		fmtParam = append(fmtParam, label)
	}

	langStr := GetInput("("+strings.Join(numbs, "|")+")", dft, "enter_co_language", strings.Join(fmtParam, "\n"))
	langNumb, _ := strconv.Atoi(langStr)

	*scriptLang = names[langNumb-1]
}

func InputForDir(dir *string, dft string, i118Key string) {
	*dir = GetInput("is_dir", dft, "enter_dir", i118Utils.Sprintf(i118Key))
}
