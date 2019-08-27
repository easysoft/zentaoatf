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
	"github.com/fatih/color"
	"os"
)

func main() {
	fmt.Fprintf(color.Output, "Windows support: %s\n", color.GreenString("PASS"))

	flagSets := make([]flag.FlagSet, 0)

	var language string
	var workDir string

	var dir string
	var langType string

	var independentFile bool
	var zentaoUrl string
	var entityType string
	var entityVal string
	var account string
	var password string

	var path string
	var files model.FlagSlice

	cuiSet := flag.NewFlagSet("atf cui - Open CUI Window", flag.ContinueOnError)
	flagSets = append(flagSets, *cuiSet)

	preferenceSet := flag.NewFlagSet("atf set - Set preferences", flag.ContinueOnError)
	flagSets = append(flagSets, *preferenceSet)
	preferenceSet.StringVar(&language, "l", "", "tool language, en or zh")
	preferenceSet.StringVar(&workDir, "d", "", "work dir")

	switchSet := flag.NewFlagSet("atf switch - Switch work dir to another path", flag.ContinueOnError)
	flagSets = append(flagSets, *switchSet)
	switchSet.StringVar(&path, "p", "", "Work dir path")

	genSet := flag.NewFlagSet("atf gen - Generate test scripts from zentaoService test cases", flag.ContinueOnError)
	flagSets = append(flagSets, *genSet)
	genSet.StringVar(&zentaoUrl, "u", "", "Zentao project url")
	genSet.StringVar(&entityType, "t", "", "Import type, 'product' or 'task'")
	genSet.StringVar(&entityVal, "v", "", "product code or task id")
	genSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	genSet.BoolVar(&independentFile, "i", false, "Save ExpectResult in independent file or not")
	genSet.StringVar(&account, "a", "", "Zentao login account")
	genSet.StringVar(&password, "p", "", "Zentao login password")

	runSet := flag.NewFlagSet("atf run - Run test scripts in specified folder", flag.ContinueOnError)
	flagSets = append(flagSets, *runSet)
	runSet.StringVar(&dir, "d", ".", "Directory that contains test scripts, base on current workdir")
	runSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	runSet.Var(&files, "f", "Script files to run, no need langType if specified, base on current workdir")

	rerunSet := flag.NewFlagSet("atf rerun - Rerun failed test scripts in specified result", flag.ContinueOnError)
	flagSets = append(flagSets, *rerunSet)
	rerunSet.StringVar(&path, "p", "", "Test result file path, base on current workdir")

	listSet := flag.NewFlagSet("atf list - List test scripts", flag.ContinueOnError)
	flagSets = append(flagSets, *listSet)
	listSet.StringVar(&dir, "d", ".", "Directory that contains test scripts")
	listSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")

	viewSet := flag.NewFlagSet("atf view - View test scripts", flag.ContinueOnError)
	flagSets = append(flagSets, *viewSet)
	viewSet.StringVar(&dir, "d", ".", "Directory that contains test scripts")
	viewSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	viewSet.Var(&files, "f", "Script files to view, no need langType if specified")

	if len(os.Args) < 2 {
		usage2(flagSets)
		os.Exit(1)
	}

	switch os.Args[1] {
	//case "mock":
	//	mock.Launch()
	case "cui":
		page.Cui()
	case "run":
		if err := runSet.Parse(os.Args[2:]); err == nil {
			if len(files) == 0 && (scriptDir == "" || langType == "") {
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
				action.GenerateScript(zentaoUrl, entityType, entityVal, langType, independentFile, account, password)
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
					action.SetLanguage(language, false)
				}
				if workDir != "" {
					action.SetWorkDir(workDir, false)
				}

				configUtils.PrintCurrConfig()
			}
		}
	default:
		usage2(flagSets)
		os.Exit(1)
	}
}

func init2() {
	if len(os.Args) > 1 {
		if os.Args[1] == "cui" {
			vari.RunFromCui = true
		} else {
			vari.RunFromCui = false
		}

		configUtils.InitConfig()
	}
}

func usage2(flagSets []flag.FlagSet) {
	fmt.Printf("Usage of atf: \n")

	for inx, flag := range flagSets {
		if inx == 0 {
			//logUtils.PrintUsageWithSpaceLine(flag, false)
		} else {
			//logUtils.PrintUsage(flag)
		}
	}

	logUtils.PrintSample()

}
