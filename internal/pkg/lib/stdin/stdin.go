package stdinUtils

import (
	"bufio"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"

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