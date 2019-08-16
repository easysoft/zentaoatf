package action

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/script"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/date"
	"strconv"
	"time"
)

func GenerateScriptFromCmd(url string, entityType string, entityVal string, langType string, singleFile bool,
	account string, password string) {
	params := make(map[string]string)

	params["entityType"] = entityType
	params["entityVal"] = entityVal

	url = commonUtils.UpdateUrl(url)
	zentaoService.Login(url, account, password)

	var productId int
	var projectId int
	var name string
	var testcases []model.TestCase
	if entityType == "product" {
		product := zentaoService.GetProductInfo(url, params["entityVal"])
		productId, _ = strconv.Atoi(product.Id)
		name = product.Name
		testcases = zentaoService.ListCaseByProduct(url, params["entityVal"])
	} else {
		task := zentaoService.GetTaskInfo(url, params["entityVal"])
		productId, _ = strconv.Atoi(task.Product)
		projectId, _ = strconv.Atoi(task.Project)
		name = task.Name
		testcases = zentaoService.ListCaseByTask(url, params["entityVal"])
	}

	if testcases != nil {
		count, err := scriptService.Generate(testcases, langType, singleFile, account, password)
		if err == nil {
			configUtils.SaveConfig("", url, params["entityType"], params["entityVal"],
				productId, projectId, langType, singleFile,
				name, account, password)

			fmt.Sprintf("success to generate %d test scripts in '%s' at %s",
				count, constant.ScriptDir, dateUtils.DateTimeStr(time.Now()))
		} else {
			fmt.Sprintf(err.Error())
		}
	}
}
