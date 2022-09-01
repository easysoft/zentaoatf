package execHelper

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	"github.com/kataras/iris/v12/websocket"
)

func Exec(ch chan int, req serverDomain.WsReq, msg *websocket.Message) (
	err error) {

	testSets := req.TestSets

	for _, testSet := range testSets {
		func() {
			serverConfig.InitExecLog(testSet.WorkspacePath)
			defer serverConfig.SyncExecLog()

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
		}() // for defer
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

func killWinProcessByUUID(uuid string) {
	cmd1 := exec.Command("cmd")
	cmd1.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c WMIC path win32_process where "CommandLine like '%%%s%%'" get Processid,Caption`, uuid), HideWindow: true}

	out, _ := cmd1.Output()
	lines := strings.Split(string(out), "\n")
	for index, line := range lines {
		if index == 0 {
			continue
		}
		line = strings.TrimSpace(line)
		cols := strings.Split(line, " ")
		if len(cols) > 3 {
			fmt.Println(fmt.Sprintf(`taskkill /F /pid %s`, cols[3]))
			cmd2 := exec.Command("cmd")
			cmd2.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c taskkill /F /pid %s`, cols[3]), HideWindow: true}
			cmd2.Start()
		}
	}
}

func killLinuxProcessByUUID(uuid string) {
	command := fmt.Sprintf(`-c ps -ef | grep %s | grep -v "grep" | awk '{print $2}' | xargs kill -9`, uuid)
	cmd := exec.Command("/bin/bash")
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: command, HideWindow: true}
	cmd.Start()

}

func KillProcessByUUID(uuid string) {
	if commonUtils.IsWin() {
		killWinProcessByUUID(uuid)
	} else {
		killLinuxProcessByUUID(uuid)
	}
}
