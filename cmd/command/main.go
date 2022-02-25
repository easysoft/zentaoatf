package main

import (
	"flag"
	"github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/command"
	"github.com/aaronchen2k/deeptest/internal/command/action"
	commandConfig "github.com/aaronchen2k/deeptest/internal/command/config"
	_consts "github.com/aaronchen2k/deeptest/internal/pkg/consts"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	resUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/res"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	"github.com/facebookgo/inject"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
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
	flagSet.StringVar(&suiteId, "suiteId", "", "")

	flagSet.StringVar(&taskId, "t", "", "")
	flagSet.StringVar(&taskId, "taskId", "", "")

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

	actionModule := injectModule()
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "run", ".")
	}

	switch os.Args[1] {
	case "set", "-set":
		action.Set()

	case "extract":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		action.Extract(files)

	case "expect":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		action.GenExpectFiles(files)

	case "checkout", "co":
		if err := flagSet.Parse(os.Args[2:]); err == nil {
			action.Generate(productId, moduleId, suiteId, taskId, independentFile, language, actionModule)
		}

	case "ci":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[len(files)+2:]); err == nil {
			action.CommitCases(files, actionModule)
		}

	case "cr":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[len(files):]); err == nil {
			action.CommitZTFTestResult(files, productId, taskId, noNeedConfirm, actionModule)
		}

	case "cb":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[len(files):]); err == nil {
			action.CommitBug(files, actionModule)
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

	case "help", "-h", "-help", "--help":
		resUtils.PrintUsage()

	case "run", "-r":
		run(os.Args, actionModule)

	default: // run
		if len(os.Args) > 1 {
			args := []string{os.Args[0], "run"}
			args = append(args, os.Args[1:]...)

			run(args, actionModule)
		} else {
			resUtils.PrintUsage()
		}
	}
}

func run(args []string, actionModule *command.IndexModule) {
	if len(args) >= 3 && stringUtils.FindInArr(args[2], _consts.UnitTestTypes) { // unit test
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

		if args[start] == _consts.UnitTestToolMvn {
			commConsts.UnitTestTool = commConsts.JUnit
			commConsts.UnitBuildTool = commConsts.Maven
		} else if args[start] == _consts.UnitTestToolRobot {
			commConsts.UnitTestTool = commConsts.RobotFramework
			commConsts.UnitBuildTool = commConsts.Maven
		}

		cmd := strings.Join(args[start:], " ")

		action.RunUnitTest(cmd)
	} else { // func test
		files := fileUtils.GetFilesFromParams(args[2:])

		err := flagSet.Parse(args[len(files)+2:])
		if err == nil {
			commConsts.ProductId = productId

			if len(files) == 0 {
				files = append(files, ".")
			}

			if commConsts.Interpreter != "" {
				msgStr := i118Utils.Sprintf("run_with_specific_interpreter", commConsts.Interpreter)
				logUtils.ExecConsolef(color.FgCyan, msgStr)
			}
			action.RunZTFTest(files, suiteId, taskId, actionModule)
		} else {
			resUtils.PrintUsage()
		}
	}
}

func init() {
	cleanup()
	commandConfig.Init()
}

func cleanup() {
	color.Unset()
}

func injectModule() (actionModule *command.IndexModule) {
	var g inject.Graph
	actionModule = command.NewIndexModule()

	// inject objects
	if err := g.Provide(
		&inject.Object{Value: dao.GetDB()},
		&inject.Object{Value: actionModule},
	); err != nil {
		logrus.Fatalf("provide usecase objects to the Graph: %v", err)
	}
	err := g.Populate()
	if err != nil {
		logrus.Fatalf("populate the incomplete Objects: %v", err)
	}
	return
}
