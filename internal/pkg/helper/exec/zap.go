package execHelper

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	"strings"
	"time"
)

type ZapScanResp struct {
	ScanId string `json:"scan"`
}
type ZapScanStatus struct {
	Status string `json:"status"`
}

var (
	OptionsMap = map[string]string{}
)

func ExecZapScan(req serverDomain.TestSet) (err error) {
	InitZap(req)

	report := commDomain.ZtfReport{
		Name:     req.Name,
		TestEnv:  commonUtils.GetOs(),
		TestType: commConsts.TestZap,
		TestTool: req.TestTool,

		StartTime: time.Now().Unix(),
	}

	report.ZapReport, _ = ExecScan()

	report.EndTime = time.Now().Unix()

	// submit result
	if req.SubmitResult && report.ZapReport != nil {
		configDir := req.WorkspacePath
		if commConsts.ExecFrom == commConsts.FromCmd {
			configDir = commConsts.ZtfDir
		}

		config := configHelper.LoadByWorkspacePath(configDir)
		err = zentaoHelper.CommitResult(report, req.ProductId, req.TaskId, config, nil)
	}

	return
}

func InitZap(req serverDomain.TestSet) {
	if req.TestTool == commConsts.Zap {
		items := strings.Split(commConsts.Options, ",")
		for _, item := range items {
			arr := strings.Split(item, "=")
			if len(arr) > 1 {
				OptionsMap[arr[0]] = arr[1]
			}
		}

		OptionsMap["server"] = httpUtils.AddSepIfNeeded(OptionsMap["server"])
		if OptionsMap["baseUrl"] == "" {
			OptionsMap["baseUrl"] = OptionsMap["site"]
		}
	}
}

func ExecScan() (report *commDomain.ZapReport, err error) {
	url := fmt.Sprintf("%sJSON/spider/action/scan/?apikey=%s&url=%s",
		httpUtils.AddSepIfNeeded(OptionsMap["server"]), OptionsMap["apiKey"], OptionsMap["site"])

	resp, err := httpUtils.Get(url)
	if err != nil {
		return
	}

	data := ZapScanResp{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return
	}

	OptionsMap["scanId"] = data.ScanId

	// wait completed
	url = fmt.Sprintf("%sJSON/spider/view/status/?apikey=%s&scanId=%s",
		OptionsMap["server"], OptionsMap["apiKey"], OptionsMap["site"])

	for i := 0; i < 100; i++ {
		resp, _ := httpUtils.Get(url)
		data := ZapScanStatus{}
		json.Unmarshal(resp, &data)

		if data.Status == "100" {
			break
		}

		time.Sleep(6 * time.Second)
	}

	url = fmt.Sprintf("%sJSON/core/view/alerts/?apikey=%s&baseurl=%s&start=0&count=10000",
		OptionsMap["server"], OptionsMap["apiKey"], OptionsMap["baseUrl"])

	resp, _ = httpUtils.Get(url)
	json.Unmarshal(resp, &report)

	return
}
