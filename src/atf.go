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

	runSet := flag.NewFlagSet("atf run: \nRun test scripts in specific folder", flag.ContinueOnError)
	runSet.StringVar(&scriptDir, "d", "", "Directory that contains test scripts")
	runSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")

	genSet := flag.NewFlagSet("atf gen: \nGenerate test scripts from zentao test cases", flag.ContinueOnError)
	genSet.StringVar(&fromUrl, "u", "", "Remote interface for test case export")
	genSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	genSet.BoolVar(&independentExpectFile, "e", false, "Save ExpectResult in an independent file or not")

	if len(os.Args) < 2 {
		fmt.Println("Usage of atf:")
		fmt.Println("  atf run -help")
		fmt.Println("  atf gen -help")
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
	}

}
