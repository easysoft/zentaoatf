package execHelper

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	"github.com/easysoft/zentaoatf/internal/server/core/dao"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	shellUtils "github.com/easysoft/zentaoatf/pkg/lib/shell"
)

func AddInterpreterIfExist(conf *commDomain.WorkspaceConf, lang string) bool {
	if commConsts.ExecFrom != commConsts.FromZentao {
		return false
	}
	data, _ := GetLangInterpreter(lang)

	if len(data) > 0 {
		var path = data[0]["path"].(string)
		configHelper.SetFieldVal(conf, lang, path)
		configHelper.UpdateAllInterpreterConfig()

		interpreter := model.Interpreter{}
		db := dao.GetDB().Model(&model.Interpreter{}).
			Where("NOT deleted").
			Where("lang = ?", lang)
		db.First(&interpreter)
		if interpreter.ID == 0 {
			interpreter = model.Interpreter{Path: path, Lang: lang}
			dao.GetDB().Model(&model.Interpreter{}).Create(&interpreter)
		}

		return true
	}

	return false
}

func GetLangInterpreter(language string) (list []map[string]interface{}, err error) {
	if commonUtils.IsWin() {
		return GetLangInterpreterWin(language)
	} else {
		return GetLangInterpreterUnix(language)
	}
}

func GetLangInterpreterUnix(language string) (list []map[string]interface{}, err error) {
	langSettings := commConsts.LangMap[language]
	whereCmd := strings.TrimSpace(langSettings["linuxWhereCmd"])
	versionCmd := strings.TrimSpace(langSettings["versionCmd"])

	output, _ := shellUtils.ExeSysCmd(whereCmd)
	pathArr := strings.Split(output, "\n")

	for _, path := range pathArr {
		path = strings.TrimSpace(path)

		if path == "" {
			continue
		}

		var vcmd string
		if language == "tcl" {
			vcmd = versionCmd + " | " + path
		} else {
			vcmd = path + " " + versionCmd + " 2>&1"
		}

		versionInfo, err1 := shellUtils.ExeSysCmd(vcmd)
		if err1 != nil {
			continue
		}

		mp := map[string]interface{}{}
		mp["path"] = path
		mp["info"] = versionInfo
		list = append(list, mp)
	}

	return
}

func GetLangInterpreterWin(language string) (list []map[string]interface{}, err error) {
	langSettings := commConsts.LangMap[language]
	whereCmd := strings.TrimSpace(langSettings["whereCmd"])
	versionCmd := strings.TrimSpace(langSettings["versionCmd"])

	path := langSettings["interpreter"]
	info := ""

	if language == "autoit" {
		if fileUtils.IsDir(filepath.Dir(path)) {
			mp := map[string]interface{}{}
			mp["path"] = path
			mp["info"] = "AutoIt V3"

			list = append(list, mp)
		}

		return
	}

	if !commonUtils.IsWin() || whereCmd == "" {
		return
	}

	output, _ := shellUtils.ExeSysCmd(whereCmd)
	pathArr := GetNoEmptyLines(strings.TrimSpace(output), ".exe", false)

	for _, path := range pathArr {
		if strings.Index(path, ".exe") != len(path)-4 {
			continue
		}
		if language == "lua" && strings.Index(path, "luac") > -1 { // compile exec file
			continue
		}

		var cmd *exec.Cmd
		if language == "tcl" {
			cmd = exec.Command("cmd", "/C", versionCmd, "|", path)
		} else {
			cmd = exec.Command("cmd", "/C", path, versionCmd)
		}

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			err = nil
			continue
		}

		infoArr := GetNoEmptyLines(out.String(), "", true)
		if len(infoArr) > 0 {
			info = infoArr[0]
		} else {
			infoArr = GetNoEmptyLines(stderr.String(), "", true)
			if len(infoArr) > 0 {
				info = infoArr[0]
			}
		}

		mp := map[string]interface{}{}
		mp["path"] = path
		mp["info"] = info
		list = append(list, mp)
	}

	return
}
func GetNoEmptyLines(text, find string, getOne bool) (ret []string) {
	arr := regexp.MustCompile("\r?\n").Split(text, -1)
	for _, item := range arr {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		if find == "" || (find != "" && strings.Contains(item, find)) {
			ret = append(ret, item)

			if getOne {
				break
			}
		}
	}

	return
}
