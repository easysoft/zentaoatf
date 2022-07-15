package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	analysisHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/analysis"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	stdinUtils "github.com/easysoft/zentaoatf/pkg/lib/stdin"
	"path/filepath"
	"strconv"
)

func CommitZTFTestResult(files []string, productId int, taskIdOrName string, noNeedConfirm bool) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		stdinUtils.InputForDir(&resultDir, "", "result")
	}

	if productId == 0 {
		productIdStr := stdinUtils.GetInput("\\d+", "",
			i118Utils.Sprintf("pls_enter")+" "+i118Utils.Sprintf("product_id"))
		productId, _ = strconv.Atoi(productIdStr)
	}

	taskName := ""
	taskId, err := strconv.Atoi(taskIdOrName)
	if err != nil {
		taskName = taskIdOrName
	}

	if taskId == 0 && taskName == "" && !noNeedConfirm {
		taskIdStr := stdinUtils.GetInput("\\d*", "",
			i118Utils.Sprintf("pls_enter")+" "+i118Utils.Sprintf("task_id")+
				i118Utils.Sprintf("task_id_empty_to_create"))
		taskId, _ = strconv.Atoi(taskIdStr)
	}

	result := serverDomain.ZentaoResultSubmitReq{
		ProductId: productId,
		TaskId:    taskId,
		Name:      taskName,
		Seq:       resultDir,
	}

	report, err := analysisHelper.ReadReportByPath(filepath.Join(result.Seq, commConsts.ResultJson))
	if err != nil {
		return
	}

	if taskName != "" {
		report.Name = taskName
	}

	config := configHelper.LoadByWorkspacePath(commConsts.ZtfDir)

	err = zentaoHelper.CommitResult(report, result.ProductId, result.TaskId, config, nil)
}
