package main

import (
	"flag"
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"os"
)

func main() {
	var exe string

	var scriptDir string
	var langType string

	var independentExpectFile bool
	var fromUrl string

	flag.StringVar(&exe, "a", "run", "action, 'run' or 'gen'")
	flag.StringVar(&scriptDir, "d", "", "Directory that contains test scripts")
	flag.StringVar(&fromUrl, "u", "", "Remote interface for test case export")
	flag.StringVar(&langType, "l", "", "Script Language like python, php etc.")
	flag.BoolVar(&independentExpectFile, "e", false, "Save ExpectResult in an independent file or not")

	flag.Parse()

	if exe == "run" {
		if langType == "" || scriptDir == "" {
			flag.Usage()
			os.Exit(1)
		} else {
			action.Run(scriptDir, langType)
		}
	} else if exe == "gen" {
		if fromUrl == "" || langType == "" {
			flag.Usage()
			os.Exit(1)
		} else {
			action.Gen(fromUrl, langType, independentExpectFile)
		}
	} else {
		fmt.Println("Usage:")

		fmt.Println("Run test scripts under specific dir")
		fmt.Println("   atf run [-d scriptDir] [-l langType]")

		fmt.Println("Generate scripts from zentao remote interface for test case export")
		fmt.Println("   atf gen [-u fromUrl] [-l langType] [-e independentExpectFile]")

		os.Exit(1)
	}

}
