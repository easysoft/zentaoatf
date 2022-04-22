package zentaoHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	websocketHelper "github.com/easysoft/zentaoatf/internal/comm/helper/websocket"
	httpUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"os"
)

func CommitResult(report commDomain.ZtfReport, productId, taskId int, config commDomain.WorkspaceConf,
	wsMsg *websocket.Message) (err error) {
	if productId != 0 {
		report.ProductId = productId
	}
	report.TaskId = taskId

	// for ci tool debug
	report.ZentaoData = os.Getenv("ZENTAO_DATA")
	report.BuildUrl = os.Getenv("BUILD_URL")

	// remove it, will cause zentao testtask not display
	//if commConsts.ComeFrom != "cmd" {
	//	report.TestType = ""
	//}

	Login(config)

	uri := fmt.Sprintf("/ciresults")
	url := GenApiUrl(uri, nil, config.Url)

	jsn, _ := json.Marshal(report)
	logUtils.Info(url)
	logUtils.Info(string(jsn))

	ret, err := httpUtils.Post(url, report)

	msg := ""
	msgColor := ""
	if err == nil {
		msg = i118Utils.Sprintf("success_to_submit_test_result")
		msgColor = color.GreenString(msg)
	} else {
		msg = i118Utils.Sprintf("fail_to_submit_test_result", err.Error())
		msgColor = color.RedString(msg)
		err = errors.New(string(ret))
	}

	logUtils.Info(msgColor)

	if commConsts.ExecFrom != commConsts.FromCmd &&
		wsMsg != nil { // from executing, not submit in webpage
		websocketHelper.SendExecMsg(msg, "", commConsts.Result, wsMsg)
	}

	return
}
