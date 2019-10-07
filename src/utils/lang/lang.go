package langUtils

import (
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var LangMap map[string]map[string]string

func initSupportedScriptLang() map[string]map[string]string {
	var once sync.Once
	once.Do(func() {
		LangMap = map[string]map[string]string{
			"bat": {
				"extName":      "bat",
				"commentsTag":  "::",
				"printGrammar": "echo #",
			},
			"javascript": {
				"extName":      "js",
				"commentsTag":  "//",
				"printGrammar": "console.log(\"#\")",
			},
			"lua": {
				"extName":      "lua",
				"commentsTag":  "--",
				"printGrammar": "print('#')",
			},
			"perl": {
				"extName":      "pl",
				"commentsTag":  "#",
				"printGrammar": "print \"#\\n\";",
			},
			"php": {
				"extName":      "php",
				"commentsTag":  "//",
				"printGrammar": "echo \"#\\n\";",
			},
			"python": {
				"extName":      "py",
				"commentsTag":  "#",
				"printGrammar": "print(\"#\")",
			},
			"ruby": {
				"extName":      "rb",
				"commentsTag":  "#",
				"printGrammar": "print(\"#\\n\")",
			},
			"shell": {
				"extName":      "sh",
				"commentsTag":  "#",
				"printGrammar": "echo \"#\"",
			},
			"tcl": {
				"extName":      "tl",
				"commentsTag":  "#",
				"printGrammar": "set hello \"#\"; \n puts [set hello];",
			},
		}
	})

	return LangMap
}

func GetSupportLanguageOptions(scriptExtsInDir []string) ([]string, []string, []string) {
	arr0 := GetSupportLanguageArrSort()

	numbs := make([]string, 0)
	names := make([]string, 0)
	labels := make([]string, 0)

	for idx, lang := range arr0 {
		ext := LangMap[lang]["extName"]

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
			lang = stringUtils.Ucfirst(lang)
		}

		labels = append(labels, strconv.Itoa(idx+1)+". "+lang)
	}

	return numbs, names, labels
}

func GetSupportLanguageArrSort() []string {
	arr := make([]string, 0)
	for lang, _ := range LangMap {
		arr = append(arr, lang)
	}

	sort.Strings(arr)

	return arr
}

func GetSupportLanguageExtArr() []string {
	arr := make([]string, 0)
	for _, key := range GetSupportLanguageArrSort() {
		arr = append(arr, LangMap[key]["extName"])
	}

	return arr
}

func CheckSupportLanguages(scriptLang string) bool {
	if LangMap[scriptLang] == nil {
		langStr := strings.Join(GetSupportLanguageArrSort(), ", ")
		logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("only_support_script_language", langStr)+"\n", color.FgRed)
		return false
	}

	return true
}

func GetSupportLanguageExtRegx() string {
	regx := "(" + strings.Join(GetSupportLanguageExtArr(), "|") + ")"

	return regx
}

func GetExtToNameMap() map[string]string {
	extMap := make(map[string]string, 0)
	for _, key := range GetSupportLanguageArrSort() {
		extMap[LangMap[key]["extName"]] = key
	}

	return extMap
}

func GetLangByFile(filePath string) string {
	ext := path.Ext(filePath)
	ext = ext[1:]

	lang := vari.ScriptExtToNameMap[ext]

	return lang
}

func init() {
	initSupportedScriptLang()
}
