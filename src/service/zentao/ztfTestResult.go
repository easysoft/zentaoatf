package zentaoService

import (
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	stdinUtils "github.com/easysoft/zentaoatf/src/utils/stdin"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"strconv"
)

func CommitZTFTestResult(resultDir string, productId string, taskId string, noNeedConfirm bool) {
	conf := configUtils.ReadCurrConfig()
	ok := Login(conf.Url, conf.Account, conf.Password)
	if !ok {
		return
	}

	report := testingService.GetZTFTestReportForSubmit(resultDir)

	if vari.ProductId == "" && productId != "" {
		vari.ProductId = productId
	}

	if taskId == "" && !noNeedConfirm {
		taskId = stdinUtils.GetInput("\\d*", "",
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("task_id")+
				i118Utils.I118Prt.Sprintf("task_id_empty_to_create"))
	}

	taskIdInt, _ := strconv.Atoi(taskId)
	CommitTestResult(report, taskIdInt)
}
