package main

import (
	"flag"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/server"
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	"github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"strconv"
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

	flagSet = flag.NewFlagSet("atf", flag.ContinueOnError)

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

	flagSet.IntVar(&vari.Port, "P", 8848, "")
	flagSet.IntVar(&vari.Port, "port", 8848, "")
	flagSet.StringVar(&vari.Platform, "M", string(serverConst.Vm), "")
	flagSet.StringVar(&vari.AgentDir, "w", "", "")

	var placeholder string
	flagSet.StringVar(&placeholder, "h", "", "")
	flagSet.StringVar(&placeholder, "r", "", "")
	flagSet.StringVar(&placeholder, "v", "", "")

	flagSet.StringVar(&vari.UnitTestResult, "result", "", "")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "run", ".")
	}

	switch os.Args[1] {
	case "run", "-r":
		run(os.Args)

	case "checkout", "co":
		if err := flagSet.Parse(os.Args[2:]); err == nil {
			action.Generate(productId, moduleId, suiteId, taskId, independentFile, language)
		}

	case "update", "up":
		if err := flagSet.Parse(os.Args[2:]); err == nil {
			action.Generate(productId, moduleId, suiteId, taskId, independentFile, language)
		}

	case "ci":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[len(files)+2:]); err == nil {
			action.CommitCases(files)
		}

	case "cr":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[len(files)+2:]); err == nil {
			action.CommitZTFTestResult(files, productId, taskId, noNeedConfirm)
		}

	case "cb":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[len(files)+2:]); err == nil {
			action.CommitBug(files)
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

	case "set", "-set":
		action.Set()

	case "sort", "-sort":
		files := fileUtils.GetFilesFromParams(os.Args[2:])
		if err := flagSet.Parse(os.Args[len(files)+2:]); err == nil {
			action.Sort(files)
		}

	case "clean", "-clean", "-c":
		action.Clean()

	case "help", "-h":
		logUtils.PrintUsage()

	default: // run
		if vari.Port != 0 {
			vari.RunMode = constant.RunModeServer
			startServer()

			return
		}

		if len(os.Args) > 1 {
			args := []string{os.Args[0], "run"}
			args = append(args, os.Args[1:]...)

			run(args)
		} else {
			logUtils.PrintUsage()
		}
	}
}

func run(args []string) {
	if len(args) >= 3 && stringUtils.FindInArr(args[2], constant.UnitTestType) { // unit test
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

		if args[start] == "mvn" {
			vari.UnitTestTool = "mvn"
		} else if args[start] == "robot" {
			vari.UnitTestTool = "robot"
		}

		cmd := strings.Join(args[start:], " ")
		action.RunUnitTest(cmd)
	} else { // func test
		files := fileUtils.GetFilesFromParams(args[2:])

		err := flagSet.Parse(args[len(files)+2:])
		if err == nil {
			vari.ProductId = productId

			if len(files) == 0 {
				files = append(files, ".")
			}

			if vari.Interpreter != "" {
				logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("run_with_specific_interpreter", vari.Interpreter), color.FgCyan)
			}
			action.RunZTFTest(files, suiteId, taskId)
		} else {
			logUtils.PrintUsage()
		}
	}
}

func startServer() {
	vari.IP = commonUtils.GetIp()
	logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("start_server", vari.IP, strconv.Itoa(vari.Port)), color.FgCyan)

	server := server.NewServer()
	server.Init()
	server.Run()

	return
}

func init() {
	cleanup()

	vari.RunFromCui = false
	configUtils.InitConfig()
}

func cleanup() {
	color.Unset()
}
