package zentaoService

import (
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	stdinUtils "github.com/easysoft/zentaoatf/src/utils/stdin"
	"strconv"
)

func CommitZTFTestResult(resultDir string, noNeedConfirm bool) {
	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	report := testingService.GetZTFTestReportForSubmit(resultDir)

	task := stdinUtils.GetInput("\\d*", "",
		i118Utils.I118Prt.Sprintf("pls_enter")+i118Utils.I118Prt.Sprintf("task_id")+
			i118Utils.I118Prt.Sprintf("task_id_empty_to_create"))

	testTask, _ := strconv.Atoi(task)
	CommitTestResult(report, testTask)
}
