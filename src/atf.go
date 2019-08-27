package main

import (
	"flag"
	"github.com/easysoft/zentaoatf/src/action"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"os"
)

func main() {
	var language string
	var independentFile bool

	//var dir string
	//var files []string
	var productId string
	var moduleId string
	var task string
	var suite string
	var result string

	flagSet := flag.NewFlagSet("atf", flag.ContinueOnError)

	//flagSet.Var(commonUtils.NewSliceValue([]string{}, &files), "f", "")
	//flagSet.Var(commonUtils.NewSliceValue([]string{}, &files), "file", "")

	//flagSet.StringVar(&dir, "d", "", "")
	//flagSet.StringVar(&dir, "dir", "", "")

	flagSet.StringVar(&productId, "p", "", "")
	flagSet.StringVar(&productId, "product", "", "")

	flagSet.StringVar(&moduleId, "m", "", "")
	flagSet.StringVar(&moduleId, "module", "", "")

	flagSet.StringVar(&suite, "s", "", "")
	flagSet.StringVar(&suite, "suite", "", "")

	flagSet.StringVar(&task, "t", "", "")
	flagSet.StringVar(&task, "task", "", "")

	flagSet.StringVar(&result, "r", "", "")
	flagSet.StringVar(&result, "result", "", "")

	flagSet.StringVar(&language, "l", "", "")
	flagSet.StringVar(&language, "language", "", "")

	flagSet.BoolVar(&independentFile, "i", false, "")
	flagSet.BoolVar(&independentFile, "independent", false, "")

	switch os.Args[1] {
	case "run":
		files, idx := commonUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[idx+1:]); err == nil {
			action.Run(files, suite, task, result)
		}

	case "checkout":
		if err := flagSet.Parse(os.Args[2:]); err == nil {
			action.GenerateScript(productId, moduleId, suite, task, independentFile, language)
		}
	case "co":
		if err := flagSet.Parse(os.Args[2:]); err == nil {
			action.GenerateScript(productId, moduleId, suite, task, independentFile, language)
		}

	case "update":
		if err := flagSet.Parse(os.Args[2:]); err == nil {
			action.GenerateScript(productId, moduleId, suite, task, independentFile, language)
		}
	case "up":
		if err := flagSet.Parse(os.Args[2:]); err == nil {
			action.GenerateScript(productId, moduleId, suite, task, independentFile, language)
		}

	case "ci":

	case "list":

	case "view":

	case "set":
		configUtils.ConfigForSet()

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
