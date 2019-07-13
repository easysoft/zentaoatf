package main

import (
	"flag"
	"github.com/easysoft/zentaoatf/src/action"
	"os"
)

func main() {
	var scriptDir string
	var langType string

	var independentExpectFile bool
	var fromUrl string

	runSet := flag.NewFlagSet("atf run", flag.ContinueOnError)
	runSet.StringVar(&scriptDir, "d", "", "Directory that contains test scripts")
	runSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")

	genSet := flag.NewFlagSet("atf gen", flag.ContinueOnError)
	genSet.StringVar(&fromUrl, "u", "", "Remote interface for test case export")
	genSet.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	genSet.BoolVar(&independentExpectFile, "e", false, "Save ExpectResult in an independent file or not")

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
