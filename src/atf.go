package main

import (
	"flag"
	"github.com/easysoft/zentaoatf/src/action"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"os"
)

func main() {
	var language string
	var independentFile bool
	var keywords string

	var productId string
	var moduleId string
	var task string
	var suite string
	var result string
	var caseId string

	flagSet := flag.NewFlagSet("atf", flag.ContinueOnError)

	//flagSet.Var(commonUtils.NewSliceValue([]string{}, &files), "file", "")

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

	flagSet.StringVar(&keywords, "k", "", "")
	flagSet.StringVar(&keywords, "keywords", "", "")

	flagSet.StringVar(&caseId, "c", "", "")
	flagSet.StringVar(&caseId, "case", "", "")

	switch os.Args[1] {
	case "run":
		files, idx := fileUtils.GetFilesFromParams(os.Args[2:])
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
		files, idx := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[idx+1:]); err == nil {
			action.CommitResult(files)
		}

	case "bug":
		files, idx := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[idx+1:]); err == nil {
			action.CommitBug(files, caseId)
		}

	case "ls":
		files, idx := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[idx+1:]); err == nil {
			action.List(files, keywords)
		}
	case "list":
		files, idx := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[idx+1:]); err == nil {
			action.List(files, keywords)
		}

	case "view":
		files, idx := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[idx+1:]); err == nil {
			action.View(files, keywords)
		}

	case "set":
		configUtils.ConfigForSet()

	case "help":
		logUtils.PrintUsage()

	default:
		logUtils.PrintUsage()
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
