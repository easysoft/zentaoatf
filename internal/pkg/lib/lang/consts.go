package langUtils

var (
	LangCommentsTagMap = map[string][]string{
		"bat":        {"goto start", ":start"},
		"javascript": {"/\\*", "\\*/"},
		"lua":        {"--\\[\\[", "\\]\\]"},
		"perl":       {"=pod", "=cut"},
		"php":        {"/\\*", "\\*/"},
		"python":     {"'''", "'''"},
		"ruby":       {"=begin", "=end"},
		"shell":      {":<<!", "!"},
		"tcl":        {"set case {", "}"},
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

	LangMap            map[string]map[string]string
	ScriptExtToNameMap map[string]string
)
