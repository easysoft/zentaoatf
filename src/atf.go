package main

import (
	"flag"
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"os"
)

func main() {
	var scriptDir string
	var langType string

	var independentExpectFile bool
	var fromUrl string

	var files strSlice

	runSet := flag.NewFlagSet("atf run: \n Run test scripts in specific folder", flag.ContinueOnError)
	runSet.StringVar(&scriptDir, "d", "./", "Directory that contains test scripts")
	runSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")

	genSet := flag.NewFlagSet("atf gen: \n Generate test scripts from zentao test cases", flag.ContinueOnError)
	genSet.StringVar(&fromUrl, "u", "", "Remote interface for test case export")
	genSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	genSet.BoolVar(&independentExpectFile, "e", false, "Save ExpectResult in an independent file or not")

	listSet := flag.NewFlagSet("atf list: \n List test scripts", flag.ContinueOnError)
	listSet.StringVar(&scriptDir, "d", "./", "Directory that contains test scripts")
	listSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")

	viewSet := flag.NewFlagSet("atf view: \n View test scripts", flag.ContinueOnError)
	viewSet.StringVar(&scriptDir, "d", "./", "Directory that contains test scripts")
	viewSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	viewSet.Var(&files, "f", "Script files to view, no need langType if specified")

	if len(os.Args) < 2 {
		fmt.Printf("Usage of atf: \n")

		fmt.Printf("atf run - Run test scripts in specific folder \n")
		runSet.PrintDefaults()

		fmt.Printf("\natf gen - Generate test scripts from zentao test cases \n")
		genSet.PrintDefaults()

		fmt.Printf("\natf list - List test cases \n")

		fmt.Printf("\natf view - View test cases \n")

		fmt.Printf("\nSample to use: \n")
		fmt.Printf("TODO... \n")

		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		if err := runSet.Parse(os.Args[2:]); err == nil {
			if langType == "" || scriptDir == "" {
				runSet.Usage()
				os.Exit(1)
			} else {
				action.Run(scriptDir, langType)
			}
		}
	case "gen":
		if err := genSet.Parse(os.Args[2:]); err == nil {
			if fromUrl == "" || langType == "" {
				genSet.Usage()
				os.Exit(1)
			} else {
				action.Gen(fromUrl, langType, independentExpectFile)
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
			if scriptDir == "" || (langType == "" || len(files) == 0) {
				viewSet.Usage()
				os.Exit(1)
			} else {
				action.View(scriptDir, langType, files)
			}
		}
	}

}

type strSlice []string

func (i *strSlice) String() string {
	return fmt.Sprintf("%d", *i)
}

func (i *strSlice) Set(value string) error {
	if value != "" {
		*i = append(*i, value)
	}
	return nil
}
