package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/command"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/stdin"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
)

func CommitZTFTestResult(files []string, productId string, taskId string, noNeedConfirm bool, actionModule *command.IndexModule) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		stdinUtils.InputForDir(&resultDir, "", "result")
	}
	if taskId == "" && !noNeedConfirm {
		taskId = stdinUtils.GetInput("\\d*", "",
			i118Utils.Sprintf("pls_enter")+" "+i118Utils.Sprintf("task_id")+
				i118Utils.Sprintf("task_id_empty_to_create"))
	}

	result := serverDomain.ZentaoResultSubmitReq{
		ProductId: productId,
		TaskId:    taskId,
		Seq:       resultDir,
	}
	actionModule.TestResultService.Submit(result, commConsts.WorkDir)
}
