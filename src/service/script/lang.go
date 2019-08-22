package scriptService

import "sync"

var LangMap map[string]map[string]string

func GetSupportedScriptLang() map[string]map[string]string {
	var once sync.Once
	once.Do(func() {
		LangMap = map[string]map[string]string{
			"bat": {
				"extName":      "bat",
				"commentsTag":  "::",
				"printGrammar": "@echo #",
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
				"printGrammar": "echo \"#\n\";",
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

func init() {
	GetSupportedScriptLang()
}
