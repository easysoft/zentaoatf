package main

import (
	"flag"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/utils/common"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"os"
)

func main() {
	var lang string

	var zentaoUrl string
	var account string
	var password string
	var language string
	var independentFile bool

	var dir string
	var files []string
	var productId string
	var moduleId string
	var taskId string
	var suite string

	flagSet := flag.NewFlagSet("atf", flag.ContinueOnError)

	flagSet.StringVar(&dir, "d", "", "")
	flagSet.Var(commonUtils.NewSliceValue([]string{}, &files), "f", "")
	flagSet.StringVar(&productId, "p", ".", "")
	flagSet.StringVar(&moduleId, "m", "", "")
	flagSet.StringVar(&taskId, "t", "", "")
	flagSet.StringVar(&suite, "s", "", "")

	switch os.Args[1] {
	case "run":

	case "co":
		if err := flagSet.Parse(os.Args[2:]); err == nil {
			action.GenerateScript(zentaoUrl, account, password,
				productId, moduleId, suite, taskId, independentFile, language)
		} else {
			color.Red(err.Error())
		}

	case "update":

	case "ci":

	case "list":

	case "view":

	case "set":
		_ = lang

	case "help":

	default:
		logUtils.PrintUsage()
		os.Exit(1)
	}
}

func init() {
	if len(os.Args) > 1 {
		if os.Args[1] == "cui" {
			vari.RunFromCui = true
		} else {
			vari.RunFromCui = false
		}
	}

	configUtils.InitConfig()
}
