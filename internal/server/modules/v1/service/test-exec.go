package service

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

var (
	ch chan int
)

type TestExecService struct {
	WorkspaceService *WorkspaceService `inject:""`
}

func NewTestExecService() *TestExecService {
	return &TestExecService{}
}

func (s *TestExecService) Start(req serverDomain.ExecReq, wsMsg *websocket.Message) (err error) {
	if execHelper.GetRunning() && req.Act != commConsts.ExecStop {
		msg := i118Utils.Sprintf("pls_stop_previous")
		websocketHelper.SendExecMsg(msg, "true", commConsts.Run, nil, wsMsg)
		logUtils.ExecConsole(color.FgRed, msg)

		return
	}

	execHelper.PopulateTestSetPropsWithParentRequest(&req)
	for idx, _ := range req.TestSets {
		testSet := &req.TestSets[idx]

		if testSet.WorkspaceId != 0 {
			po, _ := s.WorkspaceService.Get(uint(testSet.WorkspaceId))
			if testSet.WorkspacePath == "" {
				testSet.WorkspacePath = po.Path
				testSet.WorkspaceType = po.Type
			}
		}
	}

	ch = make(chan int, 1)
	go func() {
		execHelper.Exec(ch, req, wsMsg)
		execHelper.SetRunning(false)
	}()

	execHelper.SetRunning(true)

	msg := i118Utils.Sprintf("start_task")
	websocketHelper.SendExecMsg(msg, "true", commConsts.Run,
		iris.Map{"status": "start-task"}, wsMsg)
	logUtils.ExecConsole(color.FgCyan, msg)

	return
}

func (s *TestExecService) Stop(wsMsg *websocket.Message) (err error) {
	if ch != nil {
		if !execHelper.GetRunning() {
			ch = nil
		} else {
			ch <- 1
			ch = nil
		}
	}

	execHelper.SetRunning(false)

	msg := i118Utils.Sprintf("end_task")
	websocketHelper.SendExecMsg(msg, "false", commConsts.Run, nil, wsMsg)
	logUtils.ExecConsole(color.FgCyan, msg)

	return
}
