package action

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/service/script"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/date"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"time"
)

func GenerateScript(url string, entityType string, entityVal string, langType string, singleFile bool,
	account string, password string) {

	url = commonUtils.UpdateUrl(url)
	cases, productIdInt, projectId, name := zentaoService.LoadTestCases(url, account, password, entityType, entityVal)
	if cases != nil {
		count, err := scriptService.Generate(cases, langType, singleFile)
		if err == nil {
			configUtils.SaveConfig("", url, entityType, entityVal,
				productIdInt, projectId, langType, singleFile,
				name, account, password)

			configUtils.UpdateWorkDirHistoryForGenerate()

			logUtils.PrintToCmd(fmt.Sprintf("success to generate %d test scripts in '%s' at %s",
				count, constant.ScriptDir, dateUtils.DateTimeStr(time.Now())))
		} else {
			logUtils.PrintToCmd(err.Error())
		}
	}
}
