package commandConfig

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	"github.com/aaronchen2k/deeptest/internal/comm/vari"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/display"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/stdin"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"github.com/fatih/color"
	"os"
	"reflect"
	"strings"
)

type ConfigCtrl struct {
	ProjectRepo *repo.ProjectRepo `inject:""`
}

func CheckConfigPermission() {
	//err := syscall.Access(vari.ExeDir, syscall.O_RDWR)

	err := fileUtils.MkDirIfNeeded(commConsts.ExeDir + "conf")
	if err != nil {
		msg := i118Utils.Sprintf("perm_deny", commConsts.ExeDir)
		logUtils.ExecConsolef(color.FgRed, msg)
		os.Exit(0)
	}
}

func InitScreenSize() {
	w, h := display.GetScreenSize()
	vari.ScreenWidth = w
	vari.ScreenHeight = h
}

func CheckRequestConfig() {
	conf := configUtils.LoadByProjectPath(commConsts.WorkDir)
	if conf.Url == "" || conf.Username == "" || conf.Password == "" {
		InputForRequest()
	}
}

func InputForRequest() {
	conf := configUtils.LoadByProjectPath(commConsts.WorkDir)

	logUtils.ExecConsole(color.FgCyan, i118Utils.Sprintf("need_config"))

	conf.Url = stdinUtils.GetInput("(http://.*)", conf.Url, "enter_url", conf.Url)
	conf.Username = stdinUtils.GetInput("(.{2,})", conf.Username, "enter_account", conf.Username)
	conf.Password = stdinUtils.GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)

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
	if conf.Language == "zh" {
		zhCheck = "*"
		numb = "2"
	}

	numbSelected := stdinUtils.GetInput("(1|2)", numb, "enter_language", enCheck, zhCheck)

	if numbSelected == "1" {
		conf.Language = "en"
	} else {
		conf.Language = "zh"
	}

	stdinUtils.InputForBool(&configSite, true, "config_zentao_site")
	if configSite {
		conf.Url = stdinUtils.GetInput("((http|https)://.*)", conf.Url, "enter_url", conf.Url)
		conf.Url = getZenTaoBaseUrl(conf.Url)

		conf.Username = stdinUtils.GetInput("(.{2,})", conf.Username, "enter_account", conf.Username)
		conf.Password = stdinUtils.GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)
	}

	if commonUtils.IsWin() {
		var configInterpreter bool
		stdinUtils.InputForBool(&configInterpreter, true, "config_script_interpreter")
		if configInterpreter {
			scripts := scriptUtils.LoadScriptByProject(commConsts.WorkDir)
			InputForScriptInterpreter(scripts, &conf, "set")
		}
	}
	configUtils.SaveToFile(conf, commConsts.WorkDir)
	PrintCurrConfig()
}

func PrintCurrConfig() {
	logUtils.ExecConsole(color.FgCyan, "\n"+i118Utils.Sprintf("current_config"))

	val := reflect.ValueOf(vari.Config)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(vari.Config).NumField(); i++ {
		if !commonUtils.IsWin() && i > 4 {
			break
		}

		val := val.Field(i)
		name := typeOfS.Field(i).Name

		fmt.Printf("  %s: %v \n", name, val.Interface())
	}
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

		inter := stdinUtils.GetInputForScriptInterpreter(deflt, "set_script_interpreter", lang, sampleOrDefaultTips)
		configUtils.SetFieldVal(config, lang, inter)
	}

	return configChanged
}
