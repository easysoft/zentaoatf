package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	analysisUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/analysis"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/stdin"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"path/filepath"
	"strconv"
)

func CommitZTFTestResult(files []string, productId, taskId int, noNeedConfirm bool) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		stdinUtils.InputForDir(&resultDir, "", "result")
	}
	if taskId == 0 && !noNeedConfirm {
		taskIdStr := stdinUtils.GetInput("\\d*", "",
			i118Utils.Sprintf("pls_enter")+" "+i118Utils.Sprintf("task_id")+
				i118Utils.Sprintf("task_id_empty_to_create"))
		taskId, _ = strconv.Atoi(taskIdStr)
	}

	result := serverDomain.ZentaoResultSubmitReq{
		ProductId: productId,
		TaskId:    taskId,
		Seq:       resultDir,
	}

	report, err := analysisUtils.ReadReportByPath(filepath.Join(result.Seq, commConsts.ResultJson))
	if err != nil {
		return
	}

	config := configUtils.LoadByWorkspacePath(commConsts.WorkDir)
	err = zentaoUtils.CommitResult(report, result.ProductId, result.TaskId, config)
}
