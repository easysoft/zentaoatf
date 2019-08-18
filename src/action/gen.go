package action

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/service/script"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/date"
	"time"
)

func GenerateScriptFromCmd(url string, entityType string, entityVal string, langType string, singleFile bool,
	account string, password string) {

	url = commonUtils.UpdateUrl(url)
	cases, productId, projectId, name := zentaoService.LoadTestCases(url, account, password, entityType, entityVal)

	if cases != nil {
		count, err := scriptService.Generate(cases, langType, singleFile)
		if err == nil {
			configUtils.SaveConfig("", url, entityType, entityVal,
				productId, projectId, langType, singleFile,
				name, account, password)

			fmt.Sprintf("success to generate %d test scripts in '%s' at %s",
				count, constant.ScriptDir, dateUtils.DateTimeStr(time.Now()))
		} else {
			fmt.Sprintf(err.Error())
		}
	}
}
