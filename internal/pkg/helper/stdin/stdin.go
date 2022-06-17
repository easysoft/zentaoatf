package stdinHelper

import (
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stdinUtils "github.com/easysoft/zentaoatf/pkg/lib/stdin"
	"github.com/fatih/color"
	"os"
	"regexp"
	"strings"
)

func InputForScriptInterpreter(scripts []string, config *commDomain.WorkspaceConf, from string) bool {
	configChanged := false
	langs := scriptHelper.GetScriptType(scripts)

	for _, lang := range langs {
		if lang == "bat" || lang == "shell" {
			continue
		}

		deflt := configHelper.GetFieldVal(*config, lang)
		if from == "run" && deflt != "" { // already set when run, "-" means ignore
			continue
		}

		if deflt == "-" {
			deflt = ""
		}
		sampleOrDefaultTips := ""
		if deflt == "" {
			sampleOrDefaultTips = i118Utils.Sprintf("for_example", commConsts.LangMap[lang]["interpreter"]) + " " +
				i118Utils.Sprintf("empty_to_ignore")
		} else {
			sampleOrDefaultTips = deflt
		}

		configChanged = true

		inter := GetInputForScriptInterpreter(deflt, "set_script_interpreter", lang, sampleOrDefaultTips)
		configHelper.SetFieldVal(config, lang, inter)
	}

	return configChanged
}

func GetInputForScriptInterpreter(defaultVal string, fmtStr string, params ...interface{}) string {
	var ret string

	msg := i118Utils.Sprintf(fmtStr, params...)

	for {
		logUtils.ExecConsole(color.FgCyan, msg)
		stdinUtils.Scanf(&ret)

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

func InputForSet(dir string) {
	conf := configHelper.ReadFromFile(dir)

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

	numbSelected := stdinUtils.GetInput("(1|2)", numb, "enter_language", enCheck, zhCheck)

	if numbSelected == "1" {
		conf.Language = commConsts.LanguageEn
	} else {
		conf.Language = commConsts.LanguageZh
	}

	stdinUtils.InputForBool(&configSite, true, "config_zentao_site")
	if configSite {
	SetZentao:

		conf.Url = stdinUtils.GetInput("((http|https)://.*)", conf.Url, "enter_url", conf.Url)
		conf.Url = zentaoHelper.FixSiteUlt(conf.Url)
		conf.Url = fileUtils.AddUrlPathSepIfNeeded(conf.Url)

		conf.Username = stdinUtils.GetInput("(.{2,})", conf.Username, "enter_account", conf.Username)
		conf.Password = stdinUtils.GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)

		err := zentaoHelper.Login(conf)
		if err != nil {
			goto SetZentao
		}
	}

	if commonUtils.IsWin() {
		var configInterpreter bool
		stdinUtils.InputForBool(&configInterpreter, true, "config_script_interpreter")
		if configInterpreter {
			scripts := scriptHelper.LoadScriptByWorkspace(dir)
			InputForScriptInterpreter(scripts, &conf, "set")
		}
	}
	configHelper.SaveToFile(conf, dir)
}

//func getZenTaoBaseUrl(url string) string {
//	arr := strings.Split(url, "/")
//
//	base := url
//	last := arr[len(arr)-1]
//	if strings.Index(last, ".php") > -1 || strings.Index(last, ".html") > -1 ||
//		strings.Index(last, "user-login") > -1 || strings.Index(last, "?") == 0 {
//		base = base[:strings.LastIndex(base, "/")]
//	}
//
//	if strings.Index(base, "?") > -1 {
//		base = base[:strings.LastIndex(base, "?")]
//	}
//
//	return base
//}

//func InputForRequest() {
//	conf := configHelper.LoadByWorkspacePath(commConsts.WorkDir)
//
//	logUtils.ExecConsole(color.FgCyan, i118Utils.Sprintf("need_config"))
//
//	conf.Url = GetInput("(http://.*)", conf.Url, "enter_url", conf.Url)
//	conf.Username = GetInput("(.{2,})", conf.Username, "enter_account", conf.Username)
//	conf.Password = GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)
//
//	configHelper.SaveToFile(conf, commConsts.WorkDir)
//}
