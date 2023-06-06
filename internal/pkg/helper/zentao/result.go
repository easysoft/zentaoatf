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

func CommitResult(report commDomain.ZtfReport, productId, taskId int, config commDomain.WorkspaceConf,
	wsMsg *websocket.Message) (err error) {
	if productId != 0 {
		report.ProductId = productId
	}

	FilterCases(&report, config)
	report.TaskId = taskId

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

func FilterCases(report *commDomain.ZtfReport, config commDomain.WorkspaceConf) {
	//get case list
	casesResp, _ := LoadTestCaseSimple(report.ProductId, 0, 0, 0, config)

	casesMap := map[int]bool{}
	for _, caseInfo := range casesResp.Cases {
		casesMap[caseInfo.Id] = true
	}

	funcResult := make([]commDomain.FuncResult, 0)
	for _, cs := range report.FuncResult {
		if _, ok := casesMap[cs.Id]; !ok || cs.Id == 0 {
			continue
		}

		funcResult = append(funcResult, cs)
	}

	report.FuncResult = funcResult

	unitResult := make([]commDomain.UnitResult, 0)
	for _, cs := range report.UnitResult {
		cs.Id = cs.Cid

		if _, ok := casesMap[cs.Id]; !ok || cs.Id == 0 {
			continue
		}

		unitResult = append(unitResult, cs)
	}

	report.UnitResult = unitResult
	return
}
