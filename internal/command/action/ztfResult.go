package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	analysisHelper "github.com/easysoft/zentaoatf/internal/comm/helper/analysis"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/comm/helper/zentao"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	stdinUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/stdin"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
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

	report, err := analysisHelper.ReadReportByPath(filepath.Join(result.Seq, commConsts.ResultJson))
	if err != nil {
		return
	}

	config := configHelper.LoadByWorkspacePath(commConsts.WorkDir)
	err = zentaoHelper.CommitResult(report, result.ProductId, result.TaskId, config, nil)
}
