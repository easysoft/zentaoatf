package zentaoHelper

import (
	"errors"
	"fmt"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	httpUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	"github.com/fatih/color"
	"os"
)

func CommitResult(report commDomain.ZtfReport, productId, taskId int, config commDomain.WorkspaceConf) (err error) {
	if productId != 0 {
		report.ProductId = productId
	}
	report.TaskId = taskId

	// for ci tool
	report.ZentaoData = os.Getenv("ZENTAO_DATA")
	report.BuildUrl = os.Getenv("BUILD_URL")

	// remove it, will cause zentao testtask not display
	//if commConsts.ComeFrom != "cmd" {
	//	report.TestType = ""
	//}

	Login(config)

	uri := fmt.Sprintf("/ciresults")
	url := GenApiUrl(uri, nil, config.Url)

	ret, err := httpUtils.Post(url, report)

	msg := ""
	if err == nil {
		msg = color.GreenString(i118Utils.Sprintf("success_to_submit_test_result"))
	} else {
		msg = color.RedString("commit result failed, error: %s.", err.Error())
		err = errors.New(string(ret))
	}
	logUtils.Info(msg)

	return
}
