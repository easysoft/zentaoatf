package main

import (
	"flag"
	"github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/vari"
	"github.com/aaronchen2k/deeptest/internal/command"
	"github.com/aaronchen2k/deeptest/internal/command/action"
	_consts "github.com/aaronchen2k/deeptest/internal/pkg/consts"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	resUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/res"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	"github.com/fatih/color"
	"os"
	"os/signal"
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
	debug         string

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

	flagSet.StringVar(&vari.Interpreter, "interp", "", "")
	flagSet.StringVar(&vari.Interpreter, "interpreter", "", "")

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

	flagSet.StringVar(&keywords, "k", "", "")
	flagSet.StringVar(&keywords, "keywords", "", "")

	flagSet.BoolVar(&noNeedConfirm, "y", false, "")
	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	flagSet.IntVar(&vari.Port, "P", 0, "")
	flagSet.IntVar(&vari.Port, "port", 0, "")
	flagSet.StringVar(&vari.Platform, "M", string(commConsts.Vm), "")

	var placeholder string
	flagSet.StringVar(&placeholder, "h", "", "")
	flagSet.StringVar(&placeholder, "r", "", "")
	flagSet.StringVar(&placeholder, "v", "", "")

	flagSet.StringVar(&vari.UnitTestResult, "result", "", "")

	flagSet.StringVar(&debug, "debug", "", "")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "run", ".")
	}

	switch os.Args[1] {
	case "run", "-r":
		debug, os.Args = commonUtils.GetDebugParamForRun(os.Args)
		os.Setenv("debug", debug)
		//log.Println("===" + os.Getenv("debug"))
		run(os.Args)
	case "set", "-set":
		action.Set()
	case "help", "-h", "-help", "--help":
		resUtils.PrintUsage()

	default: // run
		//flagSet.Parse(os.Args[1:])
		//if vari.Port != 0 {
		//	vari.RunMode = constant.RunModeServer
		//	startServer()
		//
		//	return
		//}

		if len(os.Args) > 1 {
			args := []string{os.Args[0], "run"}
			args = append(args, os.Args[1:]...)

			run(args)
		} else {
			resUtils.PrintUsage()
		}
	}
}

func run(args []string) {
	if len(args) >= 3 && stringUtils.FindInArr(args[2], _consts.UnitTestTypes) { // unit test
		// junit -p 1 mvn clean package test
		vari.UnitTestType = args[2]
		end := 8
		if end > len(args)-1 {
			end = len(args) - 1
		}
		flagSet.Parse(args[3:])

		start := 3
		if vari.UnitTestResult != "" {
			start = start + 2
		} else {
			vari.UnitTestResult = "./"
		}
		if productId != "" {
			start = start + 2
			vari.ProductId = productId
		}
		if vari.Verbose {
			start = start + 1
		}

		if args[start] == _consts.UnitTestToolMvn {
			vari.UnitTestTool = _consts.UnitTestToolMvn
		} else if args[start] == _consts.UnitTestToolRobot {
			vari.UnitTestTool = _consts.UnitTestToolRobot
		}

		//cmd := strings.Join(args[start:], " ") todo unittest
		//
		//action.RunUnitTest(cmd)
	} else { // func test
		files := fileUtils.GetFilesFromParams(args[2:])

		err := flagSet.Parse(args[len(files)+2:])
		if err == nil {
			vari.ProductId = productId

			if len(files) == 0 {
				files = append(files, ".")
			}

			if vari.Interpreter != "" {
				logUtils.PrintToWithColor(i118Utils.Sprintf("run_with_specific_interpreter", vari.Interpreter), color.FgCyan)
			}
			action.RunZTFTest(files, suiteId, taskId)
		} else {
			resUtils.PrintUsage()
		}
	}
}

func init() {
	cleanup()
	command.InitConfig()
}

func cleanup() {
	color.Unset()
}
