package zentaoUtils

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/utils/config"
	"github.com/fatih/color"
	"os"
	"strconv"
)

func CommitResult(report commDomain.ZtfReport, productId, taskId string, projectPath string) (err error) {
	report.ProductId, _ = strconv.Atoi(productId)
	report.TaskId, _ = strconv.Atoi(taskId)

	// for ci tool
	report.ZentaoData = os.Getenv("ZENTAO_DATA")
	report.BuildUrl = os.Getenv("BUILD_URL")

	config := configUtils.LoadByProjectPath(projectPath)
	Login(config)

	url := config.Url + GenApiUri("ci", "commitResult", "")
	ret, ok := httpUtils.Post(url, report, false)

	msg := ""
	if ok {
		msg = color.GreenString(i118Utils.Sprintf("success_to_submit_test_result"))
	} else {
		msg = color.RedString(string(ret))
	}
	logUtils.Info(msg)

	return
}
