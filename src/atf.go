package main

import (
	"flag"
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
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

	var path string
	var files model.FlagSlice

	preferenceSet := flag.NewFlagSet("atf set/reset - Set preferences", flag.ContinueOnError)
	flagSets = append(flagSets, *preferenceSet)
	preferenceSet.StringVar(&language, "l", "", "tool language, en or zh")
	preferenceSet.StringVar(&workDir, "d", "./", "work dir")

	runSet := flag.NewFlagSet("atf run - Run test scripts in specified folder", flag.ContinueOnError)
	flagSets = append(flagSets, *runSet)
	runSet.StringVar(&scriptDir, "d", "./", "Directory that contains test scripts")
	runSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	runSet.Var(&files, "f", "Script files to run, no need langType if specified")

	rerunSet := flag.NewFlagSet("atf rerun - Rerun failed test scripts in specified result", flag.ContinueOnError)
	flagSets = append(flagSets, *rerunSet)
	rerunSet.StringVar(&path, "p", "", "Test result file path")

	genSet := flag.NewFlagSet("atf gen - Generate test scripts from zentao test cases", flag.ContinueOnError)
	flagSets = append(flagSets, *genSet)
	genSet.StringVar(&zentaoUrl, "u", "", "Zentao project url")
	genSet.StringVar(&entityType, "t", "", "Import from type, 'product' or 'task'")
	genSet.StringVar(&entityVal, "v", "", "product code or task id")
	genSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	genSet.BoolVar(&singleFile, "s", false, "Save ExpectResult in same file file or not")

	listSet := flag.NewFlagSet("atf list - List test scripts", flag.ContinueOnError)
	flagSets = append(flagSets, *listSet)
	listSet.StringVar(&scriptDir, "d", "./", "Directory that contains test scripts")
	listSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")

	viewSet := flag.NewFlagSet("atf view - View test scripts", flag.ContinueOnError)
	flagSets = append(flagSets, *viewSet)
	viewSet.StringVar(&scriptDir, "d", "./", "Directory that contains test scripts")
	viewSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	viewSet.Var(&files, "f", "Script files to view, no need langType if specified")

	if len(os.Args) < 2 {
		usage(flagSets)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		if err := runSet.Parse(os.Args[2:]); err == nil {
			if scriptDir == "" || (langType == "" && len(files) == 0) {
				runSet.Usage()
				os.Exit(1)
			} else {
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
	case "gen":
		if err := genSet.Parse(os.Args[2:]); err == nil {
			if zentaoUrl == "" || langType == "" || entityType == "" && entityVal == "" {
				genSet.Usage()
				os.Exit(1)
			} else {
				action.Gen(zentaoUrl, entityType, entityVal, langType, singleFile)
			}
		}
	case "list":
		if err := listSet.Parse(os.Args[2:]); err == nil {
			if scriptDir == "" || langType == "" {
				listSet.Usage()
				os.Exit(1)
			} else {
				action.List(scriptDir, langType)
			}
		}
	case "view":
		if err := viewSet.Parse(os.Args[2:]); err == nil {
			if scriptDir == "" || (langType == "" && len(files) == 0) {
				viewSet.Usage()
				os.Exit(1)
			} else {
				action.View(scriptDir, files, langType)
			}
		}

	case "set":
		if err := preferenceSet.Parse(os.Args[2:]); err == nil {
			if language == "" || workDir == "" {
				preferenceSet.Usage()
				os.Exit(1)
			} else {
				if language != "" {
					action.Set("lang", language, false)
				}
				if workDir != "" {
					action.Set("workDir", workDir, false)
				}

				utils.PrintPreference()
			}
		}
	case "reset":
		action.Reset()
	default:
		usage(flagSets)
		os.Exit(1)
	}
}

func init() {
	utils.InitPreference()
}

func usage(flagSets []flag.FlagSet) {
	fmt.Printf("Usage of atf: \n")

	for inx, flag := range flagSets {
		if inx == 0 {
			utils.PrintUsageWithSpaceLine(flag, false)
		} else {
			utils.PrintUsage(flag)
		}
	}

	utils.PrintSample()
}
