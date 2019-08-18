package main

import (
	"flag"
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/ui/page"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/config"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"os"
)

func main() {
	flagSets := make([]flag.FlagSet, 0)

	var language string
	var workDir string

	var scriptDir string
	var langType string

	var singleFile bool
	var zentaoUrl string
	var entityType string
	var entityVal string
	var account string
	var password string

	var path string
	var files model.FlagSlice

	//mockSet := flag.NewFlagSet("atf mock - Start Mock Server", flag.ContinueOnError)
	//flagSets = append(flagSets, *mockSet)

	cuiSet := flag.NewFlagSet("atf cui - Open CUI Window", flag.ContinueOnError)
	flagSets = append(flagSets, *cuiSet)

	preferenceSet := flag.NewFlagSet("atf set - Set preferences", flag.ContinueOnError)
	flagSets = append(flagSets, *preferenceSet)
	preferenceSet.StringVar(&language, "l", "", "tool language, en or zh")
	preferenceSet.StringVar(&workDir, "d", "", "work dir")

	runSet := flag.NewFlagSet("atf run - Run test scripts in specified folder", flag.ContinueOnError)
	flagSets = append(flagSets, *runSet)
	runSet.StringVar(&scriptDir, "d", ".", "Directory that contains test scripts")
	runSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	runSet.Var(&files, "f", "Script files to run, no need langType if specified")

	rerunSet := flag.NewFlagSet("atf rerun - Rerun failed test scripts in specified result", flag.ContinueOnError)
	flagSets = append(flagSets, *rerunSet)
	rerunSet.StringVar(&path, "p", "", "Test result file path")

	switchSet := flag.NewFlagSet("atf switch - Swith work dir to another path", flag.ContinueOnError)
	flagSets = append(flagSets, *switchSet)
	switchSet.StringVar(&path, "p", "", "Work dir path")

	genSet := flag.NewFlagSet("atf gen - Generate test scripts from zentaoService test cases", flag.ContinueOnError)
	flagSets = append(flagSets, *genSet)
	genSet.StringVar(&zentaoUrl, "u", "", "Zentao project url")
	genSet.StringVar(&entityType, "t", "", "Import type, 'product' or 'task'")
	genSet.StringVar(&entityVal, "v", "", "product code or task id")
	genSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	genSet.BoolVar(&singleFile, "s", false, "Save ExpectResult in same file file or not")
	genSet.StringVar(&account, "a", "", "Zentao login account")
	genSet.StringVar(&password, "p", "", "Zentao login password")

	listSet := flag.NewFlagSet("atf list - List test scripts", flag.ContinueOnError)
	flagSets = append(flagSets, *listSet)
	listSet.StringVar(&scriptDir, "d", ".", "Directory that contains test scripts")
	listSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")

	viewSet := flag.NewFlagSet("atf view - View test scripts", flag.ContinueOnError)
	flagSets = append(flagSets, *viewSet)
	viewSet.StringVar(&scriptDir, "d", ".", "Directory that contains test scripts")
	viewSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	viewSet.Var(&files, "f", "Script files to view, no need langType if specified")

	if len(os.Args) < 2 {
		usage(flagSets)
		os.Exit(1)
	}

	switch os.Args[1] {
	//case "mock":
	//	mock.Launch()
	case "cui":
		page.Cui()
	case "run":
		if err := runSet.Parse(os.Args[2:]); err == nil {
			if scriptDir == "" || (langType == "" && len(files) == 0) {
				runSet.Usage()
				os.Exit(1)
			} else {
				scriptDir = commonUtils.ConvertRunDir(scriptDir)
				action.Run(scriptDir, files, langType)
			}
		}
	case "rerun":
		if err := rerunSet.Parse(os.Args[2:]); err == nil {
			if path == "" {
				rerunSet.Usage()
				os.Exit(1)
			} else {
				action.Rerun(path)
			}
		}
	case "switch":
		if err := switchSet.Parse(os.Args[2:]); err == nil {
			if path == "" {
				switchSet.Usage()
				os.Exit(1)
			} else {
				action.SwitchWorkDir(path)
			}
		}
	case "gen":
		if err := genSet.Parse(os.Args[2:]); err == nil {
			if zentaoUrl == "" || langType == "" || entityType == "" || entityVal == "" ||
				account == "" || password == "" {
				genSet.Usage()
				os.Exit(1)
			} else {
				action.GenerateScriptFromCmd(zentaoUrl, entityType, entityVal, langType, singleFile, account, password)
			}
		}
	case "list":
		if err := listSet.Parse(os.Args[2:]); err == nil {
			if scriptDir == "" || langType == "" {
				listSet.Usage()
				os.Exit(1)
			} else {
				scriptDir = commonUtils.ConvertRunDir(scriptDir)
				action.List(scriptDir, langType)
			}
		}
	case "view":
		if err := viewSet.Parse(os.Args[2:]); err == nil {
			if scriptDir == "" || (langType == "" && len(files) == 0) {
				viewSet.Usage()
				os.Exit(1)
			} else {
				scriptDir = commonUtils.ConvertRunDir(scriptDir)
				action.View(scriptDir, files, langType)
			}
		}

	case "set":
		if err := preferenceSet.Parse(os.Args[2:]); err == nil {
			if language == "" && workDir == "" {
				preferenceSet.Usage()
				os.Exit(1)
			} else {
				if language != "" {
					action.Set("lang", language, false)
				}
				if workDir != "" {
					action.Set("workDir", workDir, false)
				}

				configUtils.PrintCurrPreference()
			}
		}
	//case "reset":
	//	action.Reset()
	default:
		usage(flagSets)
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

		configUtils.InitPreference()
	}
}

func usage(flagSets []flag.FlagSet) {
	fmt.Printf("Usage of atf: \n")

	for inx, flag := range flagSets {
		if inx == 0 {
			logUtils.PrintUsageWithSpaceLine(flag, false)
		} else {
			logUtils.PrintUsage(flag)
		}
	}

	logUtils.PrintSample()

}
