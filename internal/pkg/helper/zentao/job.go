package zentaoHelper

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/fatih/color"
)

func CommitStatus(req serverDomain.ZentaoJobSubmitReq, config commDomain.WorkspaceConf) (err error) {
	uri := fmt.Sprintf("/cijob")
	url := GenApiUrl(uri, nil, config.Url)

	jsn, _ := json.Marshal(req)
	if commConsts.Verbose {
		logUtils.Info(url)
		logUtils.Info(string(jsn))
	}

	_, err = httpUtils.Post(url, req)
	if err != nil {
		err = ZentaoRequestErr(i118Utils.Sprintf("fail_to_submit_job_status", err.Error()))
		return
	}

	msg := i118Utils.Sprintf("success_to_submit_job_status")
	logUtils.Info(color.GreenString(msg))

	return
}
