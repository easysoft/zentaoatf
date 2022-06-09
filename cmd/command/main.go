package main

import (
	"flag"
	"fmt"
	"github.com/easysoft/zentaoatf/internal/comm/consts"
	"github.com/easysoft/zentaoatf/internal/command/action"
	commandConfig "github.com/easysoft/zentaoatf/internal/command/config"
	fileUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/string"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	language        string
	independentFile bool
	keywords        string

	productId string
	moduleId  string
	taskId    string
	suiteId   string

	noNeedConfirm bool

	flagSet *flag.FlagSet
)

func main() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		cleanup()
		os.Exit(0)
	}()

	flagSet = flag.NewFlagSet("ztf", flag.ContinueOnError)

	flagSet.StringVar(&productId, "p", "", "")
	flagSet.StringVar(&productId, "product", "", "")

	flagSet.StringVar(&moduleId, "m", "", "")
	flagSet.StringVar(&moduleId, "module", "", "")

	flagSet.StringVar(&suiteId, "s", "", "")
	flagSet.StringVar(&suiteId, "suite", "", "")

	flagSet.StringVar(&taskId, "t", "", "")
	flagSet.StringVar(&taskId, "task", "", "")

	flagSet.StringVar(&language, "l", "", "")
	flagSet.StringVar(&language, "language", "", "")

	flagSet.BoolVar(&independentFile, "i", false, "")
	flagSet.BoolVar(&independentFile, "independent", false, "")

	flagSet.StringVar(&commConsts.Interpreter, "I", "", "")
	flagSet.StringVar(&commConsts.Interpreter, "interpreter", "", "")

	flagSet.StringVar(&keywords, "k", "", "")
	flagSet.StringVar(&keywords, "keywords", "", "")

	flagSet.BoolVar(&noNeedConfirm, "y", false, "")
	flagSet.BoolVar(&commConsts.Verbose, "verbose", false, "")

	flagSet.StringVar(&commConsts.UnitTestResult, "result", "", "")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "run", ".")
	}

	switch os.Args[1] {
	case "set", "-set":
		flagSet.Parse(os.Args[2:])
		if commConsts.Verbose {
			fmt.Printf("\nIsRelease=%t\n", commConsts.IsRelease)
			fmt.Printf("\nlaunch %s%s in %s\n", "", commConsts.App, commConsts.WorkDir)
		}
		action.Set()

	case "expect":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		action.GenExpectFiles(files)
	case "extract":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		action.Extract(files)

	case "checkout", "co":
		if err := flagSet.Parse(os.Args[2:]); err == nil {
			action.Checkout(productId, moduleId, suiteId, taskId, independentFile, language)
		}

	case "ci":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[len(files)+2:]); err == nil {
			action.CommitCases(files)
		}

	case "cr":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[len(files)+2:]); err == nil {
			action.CommitZTFTestResult(files, stringUtils.ParseInt(productId), stringUtils.ParseInt(taskId),
				noNeedConfirm)
		}

	case "cb":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[len(files)+2:]); err == nil {
			action.CommitBug(files, stringUtils.ParseInt(productId))
		}

	case "list", "ls", "-l":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[len(files)+2:]); err == nil {
			if len(files) == 0 {
				files = append(files, ".")
			}

			action.List(files, keywords)
		}

	case "view", "-v":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[len(files)+2:]); err == nil {
			action.View(files, keywords)
		}

	case "clean", "-clean", "-c":
		action.Clean()

	case "version", "--version":
		logUtils.PrintVersion(commConsts.AppVersion, commConsts.BuildTime, commConsts.GoVersion, commConsts.GitHash)

	case "help", "-h", "-help", "--help":
		action.PrintUsage()

	case "run", "-r":
		run(os.Args)

	default: // run
		if len(os.Args) > 1 {
			args := []string{os.Args[0], "run"}
			args = append(args, os.Args[1:]...)

			run(args)
		} else {
			action.PrintUsage()
		}
	}
}

func run(args []string) {
	if len(args) >= 3 && stringUtils.FindInArr(args[2], commConsts.UnitTestTypes) { // unit test
		runUnitTest(args)
	} else { // ztf test
		runFuncTest(args)
	}
}

func runFuncTest(args []string) {
	files := fileUtils.GetFilesFromParams(args[2:])

	err := flagSet.Parse(args[len(files)+2:])
	if err != nil {
		action.PrintUsage()
		return
	}
	if len(files) > 0 && files[0] != "" && !fileUtils.FileExist(files[0]) {
		action.PrintUsage()
		return
	}

	commConsts.ProductId = productId

	if len(files) == 0 {
		files = append(files, ".")
	}

	if commConsts.Interpreter != "" {
		msgStr := i118Utils.Sprintf("run_with_specific_interpreter", commConsts.Interpreter)
		logUtils.ExecConsolef(color.FgCyan, msgStr)
	}
	action.RunZTFTest(files, moduleId, suiteId, taskId)
}

func runUnitTest(args []string) {
	// junit -p 1 mvn clean package test
	commConsts.UnitTestType = args[2]
	end := 8
	if end > len(args)-1 {
		end = len(args) - 1
	}
	flagSet.Parse(args[3:])

	start := 3
	if commConsts.UnitTestResult != "" {
		start = start + 2
	} else {
		commConsts.UnitTestResult = "./"
	}
	if productId != "" {
		start = start + 2
		commConsts.ProductId = productId
	}
	if commConsts.Verbose {
		start = start + 1
	}

	if args[start] == commConsts.UnitTestToolMvn {
		commConsts.UnitTestTool = commConsts.JUnit
		commConsts.UnitBuildTool = commConsts.Maven
	} else if args[start] == commConsts.UnitTestToolRobot {
		commConsts.UnitTestTool = commConsts.RobotFramework
		commConsts.UnitBuildTool = commConsts.Maven
	}

	cmd := strings.Join(args[start:], " ")

	action.RunUnitTest(cmd)
}

func init() {
	cleanup()
	commandConfig.Init()
}

func cleanup() {
	color.Unset()
}
