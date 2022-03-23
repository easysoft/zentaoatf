package langHelper

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	"path"
	"sort"
	"strconv"
	"strings"
)

func GetSupportLanguageOptions(scriptExtsInDir []string) ([]string, []string, []string) {
	arr0 := GetSupportLanguageArrSort()

	numbs := make([]string, 0)
	names := make([]string, 0)
	labels := make([]string, 0)

	for idx, lang := range arr0 {
		ext := commConsts.LangMap[lang]["extName"]

		if scriptExtsInDir != nil {
			found := stringUtils.FindInArr(ext, scriptExtsInDir)
			if !found {
				continue
			}
		}

		numbs = append(numbs, strconv.Itoa(idx+1))
		names = append(names, lang)

		if lang == "bat" || lang == "php" {
			lang = stringUtils.UcAll(lang)
		} else {
			lang = stringUtils.UcFirst(lang)
		}

		labels = append(labels, strconv.Itoa(idx+1)+". "+lang)
	}

	return numbs, names, labels
}

func GetSupportLanguageArrSort() []string {
	arr := make([]string, 0)
	for lang, _ := range commConsts.LangMap {
		arr = append(arr, lang)
	}

	sort.Strings(arr)

	return arr
}

func GetSupportLanguageExtArr() []string {
	arr := make([]string, 0)
	for _, key := range GetSupportLanguageArrSort() {
		arr = append(arr, commConsts.LangMap[key]["extName"])
	}

	return arr
}

func CheckSupportLanguages(scriptLang string) bool {
	if commConsts.LangMap[scriptLang] == nil {
		langStr := strings.Join(GetSupportLanguageArrSort(), ", ")
		logUtils.Errorf(i118Utils.Sprintf("only_support_script_language", langStr))
		return false
	}

	return true
}

func GetSupportLanguageExtRegx() string {
	regx := "(" + strings.Join(GetSupportLanguageExtArr(), "|") + ")"

	return regx
}

func GetExtToNameMap() {
	if commConsts.ScriptExtToNameMap != nil { // init once
		return
	}

	commConsts.ScriptExtToNameMap = make(map[string]string, 0)
	for _, key := range GetSupportLanguageArrSort() {
		commConsts.ScriptExtToNameMap[commConsts.LangMap[key]["extName"]] = key
	}

	return
}

func GetLangByFile(filePath string) string {
	ext := path.Ext(filePath)
	ext = ext[1:]

	lang := commConsts.ScriptExtToNameMap[ext]

	return lang
}
