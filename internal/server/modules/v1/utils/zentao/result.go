package zentaoUtils

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/config"
	"github.com/fatih/color"
	"os"
	"strconv"
)

func CommitResult(report commDomain.ZtfReport, result serverDomain.ZentaoResultSubmitReq, projectPath string) (err error) {
	report.ProductId, _ = strconv.Atoi(result.ProductId)
	report.TaskId, _ = strconv.Atoi(result.TaskId)

	// for ci tool
	report.ZentaoData = os.Getenv("ZENTAO_DATA")
	report.BuildUrl = os.Getenv("BUILD_URL")

	config := configUtils.LoadByProjectPath(projectPath)
	Login(config)

	url := config.Url + GenApiUri("ci", "commitResult", "")
	bytes, ok := httpUtils.Post(url, report, false)
	if ok {
		err = GetRespErr(bytes, "fail_to_submit_test_result")
	}

	msg := ""
	if err == nil {
		msg = color.GreenString(i118Utils.Sprintf("success_to_submit_test_result"))
	} else {
		msg = color.RedString(err.Error())
	}
	logUtils.Info(msg)

	return
}
