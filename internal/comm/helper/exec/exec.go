package execHelper

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12/websocket"
	"strings"
)

func Exec(ch chan int, req serverDomain.WsReq, msg *websocket.Message) (
	err error) {

	testSets := req.TestSets

	for _, testSet := range testSets {
		serverConfig.InitExecLog(testSet.WorkspacePath)

		if testSet.ScriptDirParamFromCmdLine == "" {
			testSet.ScriptDirParamFromCmdLine = req.ScriptDirParamFromCmdLine
		}

		if req.Act == commConsts.ExecCase {
			ExecCases(ch, testSet, msg)
		} else if req.Act == commConsts.ExecModule {
			ExecModule(ch, testSet, msg)
		} else if req.Act == commConsts.ExecSuite {
			ExecSuite(ch, testSet, msg)
		} else if req.Act == commConsts.ExecTask {
			ExecTask(ch, testSet, msg)
		} else if req.Act == commConsts.ExecUnit {
			ExecUnit(ch, testSet, msg)
		}

		logUtils.LoggerExecFile = nil
		logUtils.LoggerExecResult = nil
	}

	return
}

func PopulateTestSetProps(req *serverDomain.WsReq) {
	for idx, _ := range req.TestSets {
		testSet := &req.TestSets[idx]

		testSet.Scope = req.Scope
		testSet.Seq = req.Seq

		if testSet.ProductId == 0 && req.ProductId != 0 {
			testSet.ProductId = req.ProductId
		}

		testSet.ModuleId = req.ModuleId
		testSet.SuiteId = req.SuiteId
		testSet.TaskId = req.TaskId

		testSet.ScriptDirParamFromCmdLine = req.ScriptDirParamFromCmdLine

		setTestTool(testSet, *req)
		setBuildTool(testSet, *req)

		if testSet.Cmd == "" && req.Cmd != "" {
			testSet.Cmd = req.Cmd
		}

		if !testSet.SubmitResult && req.SubmitResult {
			testSet.SubmitResult = req.SubmitResult
		}

	}
}

func setTestTool(testSet *serverDomain.TestSet, req serverDomain.WsReq) {
	if testSet.TestTool == "" && req.TestTool != "" {
		testSet.TestTool = req.TestTool
	}

	if testSet.TestTool == "" {
		testSet.TestTool = testSet.WorkspaceType
	}
}

func setBuildTool(testSet *serverDomain.TestSet, req serverDomain.WsReq) {
	if testSet.BuildTool == "" && req.BuildTool != "" {
		testSet.BuildTool = req.BuildTool
	}

	if testSet.BuildTool == "" {
		arr := strings.Split(testSet.Cmd, " ")
		testSet.BuildTool = commConsts.UnitBuildToolMap[strings.TrimSpace(arr[0])]
	}
}
