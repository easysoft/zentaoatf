package zentaoHelper

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/bitly/go-simplejson"
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

	//convert steps array to map
	newReport := convertStepsToMap(jsn)

	if commConsts.Verbose {
		logUtils.Info(url)
		logUtils.Info(string(jsn))
	}

	_, err = httpUtils.Post(url, newReport)
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

	ignoredCases := map[int]bool{}

	funcResult := make([]commDomain.FuncResult, 0)
	for _, cs := range report.FuncResult {
		if _, ok := casesMap[cs.Id]; !ok && cs.Id != 0 {
			ignoredCases[cs.Id] = true
			continue
		}

		funcResult = append(funcResult, cs)
	}

	report.FuncResult = funcResult

	unitResult := make([]commDomain.UnitResult, 0)
	for _, cs := range report.UnitResult {
		cs.Id = cs.Cid

		if _, ok := casesMap[cs.Id]; !ok || cs.Id == 0 {
			ignoredCases[cs.Id] = true
			continue
		}

		unitResult = append(unitResult, cs)
	}

	report.UnitResult = unitResult

	if len(ignoredCases) > 0 {
		msg := i118Utils.Sprintf("ignored_cases_to_submit_test_result", report.ProductId)
		logUtils.Info("\n" + color.HiYellowString(msg))
		for k := range ignoredCases {
			logUtils.Info("  " + strconv.Itoa(k))
		}
	}
	return
}

func convertStepsToMap(rawJsn []byte) interface{} {
	jsn, err := simplejson.NewJson(rawJsn)
	if err != nil {
		return nil
	}
	funcResult, err := jsn.Get("funcResult").Array()
	if err != nil || funcResult == nil {
		return jsn.Interface()
	}
	for index, row := range funcResult {
		result := row.(map[string]interface{})
		steps := result["steps"].([]interface{})
		newSteps := map[string]interface{}{}

		for _, step := range steps {
			stepMap := step.(map[string]interface{})
			newSteps[stepMap["id"].(string)] = step
		}

		result["steps"] = newSteps

		funcResult[index] = result
	}

	jsn.Set("funcResult", funcResult)

	return jsn.Interface()
}
