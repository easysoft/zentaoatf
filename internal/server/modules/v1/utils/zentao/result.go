package zentaoUtils

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/config"
	"github.com/bitly/go-simplejson"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
)

func CommitResult(report commDomain.ZtfReport, result serverDomain.ZentaoResult, projectPath string) {
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
		json, err1 := simplejson.NewJson(bytes)
		if err1 == nil {
			result, err2 := json.Get("result").String()
			if err2 != nil || result != "success" {
				ok = false
			}
		} else {
			ok = false
		}
	}

	msg := ""
	if ok {
		msg += color.GreenString(i118Utils.Sprintf("success_to_submit_test_result"))
	} else {
		msg = i118Utils.Sprintf("fail_to_submit_test_result")
		if strings.Index(string(bytes), "login") > -1 {
			msg = i118Utils.Sprintf("fail_to_login")
		}
		msg = color.RedString(msg)
	}
	logUtils.Infof(msg)
}
