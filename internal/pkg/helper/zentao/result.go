package zentaoHelper

import (
	"encoding/json"
	"fmt"
	"os"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
)

func CommitResult(report commDomain.ZtfReport, productId, taskId, task int, config commDomain.WorkspaceConf,
	wsMsg *websocket.Message) (err error) {
	if productId != 0 {
		report.ProductId = productId
	}
	RemoveAutoCreateId(&report)
	report.TaskId = taskId
	report.Task = task

	// for ci tool debug
	report.ZentaoData = os.Getenv("ZENTAO_DATA")
	report.BuildUrl = os.Getenv("BUILD_URL")

	if commConsts.ExecFrom != commConsts.FromZentao {
		Login(config)
	}

	uri := fmt.Sprintf("/ciresults")
	url := GenApiUrl(uri, nil, config.Url)

	jsn, _ := json.Marshal(report)
	if commConsts.Verbose {
		logUtils.Info(url)
		logUtils.Info(string(jsn))
	}

	_, err = httpUtils.Post(url, report)
	if err != nil {
		err = ZentaoRequestErr(i118Utils.Sprintf("fail_to_submit_test_result", err.Error()))
		return
	}

	msg := i118Utils.Sprintf("success_to_submit_test_result")
	logUtils.Info(color.GreenString(msg))

	if commConsts.ExecFrom == commConsts.FromClient &&
		wsMsg != nil { // from executing, not submit in webpage
		websocketHelper.SendExecMsg(msg, "", commConsts.Result, nil, wsMsg)
	}

	return
}

func JobCommitResult(report interface{}, config commDomain.WorkspaceConf) (err error) {
	uri := "/ztf/submitResult"
	url := GenApiUrl(uri, nil, config.Url)

	if commConsts.Verbose {
		jsn, _ := json.Marshal(report)
		logUtils.Info(url)
		logUtils.Info(string(jsn))
	}

	_, err = httpUtils.Post(url, report)
	if err != nil {
		err = ZentaoRequestErr(i118Utils.Sprintf("fail_to_submit_test_result", err.Error()))
		return
	}

	msg := i118Utils.Sprintf("success_to_submit_test_result")
	logUtils.Info(color.GreenString(msg))

	return
}

func RemoveAutoCreateId(report *commDomain.ZtfReport) {
	if report.TestType == commConsts.TestFunc {
		return
	}
	for idx, cs := range report.UnitResult {
		report.UnitResult[idx].Id = cs.Cid
	}
	return
}
