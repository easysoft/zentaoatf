package zentaoService

import (
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/fatih/color"
)

func CommitUnitTestResult(report model.UnitTestReport) {
	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	url := conf.Url + zentaoUtils.GenApiUri("unittest", "commitResult", "")
	_, ok := client.PostObject(url, report)

	msg := ""
	if ok {
		msg = color.GreenString(i118Utils.I118Prt.Sprintf("success_to_submit_unit_test_result"))
	} else {
		msg = color.RedString(i118Utils.I118Prt.Sprintf("fail_to_submit_unit_test_result"))
	}
	logUtils.Screen(msg)
}
