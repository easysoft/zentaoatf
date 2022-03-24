package commConsts

var (
	ScriptExtToNameMap map[string]string
)

var (
	LangMap = map[string]map[string]string{
		"shell": {
			"name":         "Shell",
			"extName":      "sh",
			"commentsTag":  "#",
			"printGrammar": "echo \"#\"",
		},
		"bat": {
			"name":         "BAT",
			"extName":      "bat",
			"commentsTag":  "::",
			"printGrammar": "echo #",
		},
		"autoit": {
			"name":         "AutoIT",
			"extName":      "au3",
			"commentsTag":  "#",
			"printGrammar": "ConsoleWrite(text & @CRLF)",
			"interpreter":  "c:\\Program Files (x86)\\AutoIt3\\AutoIt3_x64.exe",
		},

		"javascript": {
			"name":         "JavaScript",
			"extName":      "js",
			"commentsTag":  "//",
			"printGrammar": "console.log(\"#\")",
			"interpreter":  "C:\\Program Files\\nodejs\\node.exe",
			"whereCmd":     "where node",
			"versionCmd":   "-v",
		},
		"lua": {
			"name":         "Lua",
			"extName":      "lua",
			"commentsTag":  "--",
			"printGrammar": "print('#')",
			"interpreter":  "C:\\lua-5.3.6_Win32_bin",
			"whereCmd":     "where lua*.exe",
			"versionCmd":   "-v",
		},
		"perl": {
			"name":         "Perl",
			"extName":      "pl",
			"commentsTag":  "#",
			"printGrammar": "print \"#\\n\";",
			"interpreter":  "C:\\Perl64\\bin\\perl.exe",
			"whereCmd":     "where perl",
			"versionCmd":   "-v",
		},
		"php": {
			"name":         "PHP",
			"extName":      "php",
			"commentsTag":  "//",
			"printGrammar": "echo \"#\\n\";",
			"interpreter":  "C:\\php-7.3.9-Win32-VC15-x64\\php.exe",
			"whereCmd":     "where php",
			"versionCmd":   "-v",
		},
		"python": {
			"name":         "Python",
			"extName":      "py",
			"commentsTag":  "#",
			"printGrammar": "print(\"#\")",
			"interpreter":  "C:\\Users\\admin\\AppData\\Local\\Programs\\Python\\Python37-32\\python.exe",
			"whereCmd":     "where python",
			"versionCmd":   "--version",
		},
		"ruby": {
			"name":         "Ruby",
			"extName":      "rb",
			"commentsTag":  "#",
			"printGrammar": "print(\"#\\n\")",
			"interpreter":  "C:\\Ruby26-x64\\bin\\ruby.exe",
			"whereCmd":     "where ruby",
			"versionCmd":   "-v",
		},
		"tcl": {
			"name":         "TCL",
			"extName":      "tl",
			"commentsTag":  "#",
			"printGrammar": "set hello \"#\"; \n puts [set hello];",
			"interpreter":  "C:\\ActiveTcl\\bin\\tclsh.exe",
			"whereCmd":     "where tclsh",
			"versionCmd":   "echo puts $tcl_version;exit 0",
		},
	}

	LangCommentsTagMap = map[string][]string{
		"bat":        {`goto start`, `:start`},
		"javascript": {`/*`, `*/`},
		"lua":        {`--[[`, `]]`},
		"perl":       {`=pod`, `=cut`},
		"php":        {`/**`, `*/`},
		"python":     {"'''", "'''"},
		"ruby":       {`=begin`, `=end`},
		"shell":      {`:<<!`, `!`},
		"tcl":        {`set case {`, `}`},
	}

	LangCommentsRegxMap = map[string][]string{
		"bat":        {"^\\s*goto start\\s*$", "^\\s*:start\\s*$"},
		"javascript": {"^\\s*/\\*{1,}\\s*$", "^\\s*\\*{1,}/\\s*$"},
		"lua":        {"^\\s*--\\[\\[\\s*$", "^\\s*\\]\\]\\s*$"},
		"perl":       {"^\\s*=pod\\s*$", "^\\s*=cut\\s*$"},
		"php":        {"^\\s*/\\*{1,}\\s*$", "^\\s*\\*{1,}/\\s*$"},
		"python":     {"^\\s*'''\\s*$", "^\\s*'''\\s*$"},
		"ruby":       {"^\\s*=begin\\s*$", "^\\s*=end\\s*$"},
		"shell":      {"^\\s*:<<!\\s*$", "^\\s*!\\s*$"},
		"tcl":        {"^\\s*set case {", "^\\s*}"},
	}
)
