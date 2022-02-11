package stdinUtils

import (
	"bufio"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
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
		logUtils.Info(msg)
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

			logUtils.Info(ret)
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
		} else {
			pass, _ = regexp.MatchString("^"+regx+"$", temp)
			msg = "invalid_input"
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

func GetInputForScriptInterpreter(defaultVal string, fmtStr string, params ...interface{}) string {
	var ret string

	msg := i118Utils.Sprintf(fmtStr, params...)

	for {
		logUtils.ExecConsole(color.FgCyan, msg)
		Scanf(&ret)

		ret = strings.TrimSpace(ret)

		if ret == "" && defaultVal != "" {
			ret = defaultVal

			logUtils.Info(ret)
		}

		if ret == "exit" {
			color.Unset()
			os.Exit(0)
		}

		if ret == "" { // ignore to set
			return "-"
		}

		sep := string(os.PathSeparator)
		if sep == `\` {
			sep = `\\`
		}
		reg := fmt.Sprintf(".*%s+[^%s]+", sep, sep)
		pass, _ := regexp.MatchString(reg, ret)
		if pass {
			return ret
		} else {
			ret = ""
			logUtils.ExecConsole(color.FgRed, i118Utils.Sprintf("invalid_input"))
		}
	}
}

func InputForCheckout(productId *string, moduleId *string, suiteId *string, taskId *string,
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

	numbs, names, labels := langUtils.GetSupportLanguageOptions(nil)
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

func InputForRequest() {
	conf := configUtils.LoadByProjectPath(commConsts.WorkDir)

	logUtils.ExecConsole(color.FgCyan, i118Utils.Sprintf("need_config"))

	conf.Url = GetInput("(http://.*)", conf.Url, "enter_url", conf.Url)
	conf.Username = GetInput("(.{2,})", conf.Username, "enter_account", conf.Username)
	conf.Password = GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)

	configUtils.SaveToFile(conf, commConsts.WorkDir)
}

func InputForSet() {
	conf := configUtils.ReadFromFile(commConsts.WorkDir)

	var configSite bool

	logUtils.ExecConsole(color.FgCyan, i118Utils.Sprintf("begin_config"))

	enCheck := ""
	var numb string
	if conf.Language == "en" {
		enCheck = "*"
		numb = "1"
	}
	zhCheck := ""
	if conf.Language == commConsts.LanguageZh {
		zhCheck = "*"
		numb = "2"
	}

	numbSelected := GetInput("(1|2)", numb, "enter_language", enCheck, zhCheck)

	if numbSelected == "1" {
		conf.Language = commConsts.LanguageEn
	} else {
		conf.Language = commConsts.LanguageZh
	}

	InputForBool(&configSite, true, "config_zentao_site")
	if configSite {
		conf.Url = GetInput("((http|https)://.*)", conf.Url, "enter_url", conf.Url)
		conf.Url = getZenTaoBaseUrl(conf.Url)

		conf.Username = GetInput("(.{2,})", conf.Username, "enter_account", conf.Username)
		conf.Password = GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)
	}

	if commonUtils.IsWin() {
		var configInterpreter bool
		InputForBool(&configInterpreter, true, "config_script_interpreter")
		if configInterpreter {
			scripts := scriptUtils.LoadScriptByProject(commConsts.WorkDir)
			InputForScriptInterpreter(scripts, &conf, "set")
		}
	}
	configUtils.SaveToFile(conf, commConsts.WorkDir)
}

func getZenTaoBaseUrl(url string) string {
	arr := strings.Split(url, "/")

	base := url
	last := arr[len(arr)-1]
	if strings.Index(last, ".php") > -1 || strings.Index(last, ".html") > -1 ||
		strings.Index(last, "user-login") > -1 || strings.Index(last, "?") == 0 {
		base = base[:strings.LastIndex(base, "/")]
	}

	if strings.Index(base, "?") > -1 {
		base = base[:strings.LastIndex(base, "?")]
	}

	return base
}

func InputForScriptInterpreter(scripts []string, config *commDomain.ProjectConf, from string) bool {
	configChanged := false
	langs := scriptUtils.GetScriptType(scripts)

	for _, lang := range langs {
		if lang == "bat" || lang == "shell" {
			continue
		}

		deflt := configUtils.GetFieldVal(*config, lang)
		if from == "run" && deflt != "" { // already set when run, "-" means ignore
			continue
		}

		if deflt == "-" {
			deflt = ""
		}
		sampleOrDefaultTips := ""
		if deflt == "" {
			sampleOrDefaultTips = i118Utils.Sprintf("for_example", langUtils.LangMap[lang]["interpreter"]) + " " +
				i118Utils.Sprintf("empty_to_ignore")
		} else {
			sampleOrDefaultTips = deflt
		}

		configChanged = true

		inter := GetInputForScriptInterpreter(deflt, "set_script_interpreter", lang, sampleOrDefaultTips)
		configUtils.SetFieldVal(config, lang, inter)
	}

	return configChanged
}

func InputForDir(dir *string, dft string, i118Key string) {
	*dir = GetInput("is_dir", dft, "enter_dir", i118Utils.Sprintf(i118Key))
}
