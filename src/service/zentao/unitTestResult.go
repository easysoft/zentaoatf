package zentaoService

import (
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/fatih/color"
	"os"
)

func CommitUnitTestResult(report model.UnitTestReport) {
	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	buildNumb := os.Getenv("BUILD_NUMBER")
	url := conf.Url + zentaoUtils.GenApiUri("unittest", "commitResult", buildNumb)
	//logUtils.Screen(url)
	resp, ok := client.PostObject(url, report)

	if ok {
		json, err1 := simplejson.NewJson([]byte(resp))
		if err1 == nil {
			result, err2 := json.Get("result").String()
			if err2 != nil || result != "success" {
				ok = false
			}
		} else {
			ok = false
		}
	}

	msg := "\n"
	if ok {
		msg += color.GreenString(i118Utils.I118Prt.Sprintf("success_to_submit_unit_test_result"))
	} else {
		msg += color.RedString(i118Utils.I118Prt.Sprintf("fail_to_submit_unit_test_result", url))
		msg += "\n" + i118Utils.I118Prt.Sprintf("server_return")
		msg += "\n" + resp
	}
	logUtils.Screen(msg)
}
