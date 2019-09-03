package langUtils

import (
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/fatih/color"
	"strconv"
	"strings"
	"sync"
)

var LangMap map[string]map[string]string

func GetSupportedScriptLang() map[string]map[string]string {
	var once sync.Once
	once.Do(func() {
		LangMap = map[string]map[string]string{
			"bat": {
				"extName":      "bat",
				"commentsTag":  "::",
				"printGrammar": "echo #",
			},
			"go": {
				"extName":      "go",
				"commentsTag":  "//",
				"printGrammar": "println(\"#\")",
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

func GetSupportLanguageOptions() ([]string, []string, []string) {
	arr0 := GetSupportLangageArr()

	numbs := make([]string, 0)
	names := make([]string, 0)
	labels := make([]string, 0)

	for idx, lang := range arr0 {
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

func GetSupportLangageArr() []string {
	langMap := GetSupportedScriptLang()

	arr := make([]string, 0)
	for lang, _ := range langMap {
		arr = append(arr, lang)
	}

	return arr
}

func CheckSupportLangages(scriptLang string) bool {
	langMap := GetSupportedScriptLang()

	if langMap[scriptLang] == nil {
		langs := strings.Join(GetSupportLangageArr(), ", ")
		logUtils.PrintToCmd(i118Utils.I118Prt.Sprintf("only_support_script_language", langs)+"\n", color.FgRed)
		return false
	}

	return true
}

func GetSupportLangageRegx() string {
	regx := "(" + strings.Join(GetSupportLangageArr(), "|") + ")"

	return regx
}

func init() {
	GetSupportedScriptLang()
}
