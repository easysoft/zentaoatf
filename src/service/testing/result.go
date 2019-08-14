package testingService

import (
	"encoding/json"
	httpClient "github.com/easysoft/zentaoatf/src/service/client"
	"github.com/easysoft/zentaoatf/src/utils/common"
	config2 "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	print2 "github.com/easysoft/zentaoatf/src/utils/print"
	"path"
	"strconv"
	"strings"
)

func SubmitResult(caseList []string) {
	config := config2.ReadCurrConfig()

	entityType := config.EntityType
	entityVal := config.EntityVal

	requestObj := make(map[string]string)
	requestObj["entityType"] = entityType
	requestObj["entityVal"] = entityVal

	cases := make(map[int]bool)
	for _, str := range caseList {
		arr := strings.Split(str, " ")
		var status bool
		str := strings.ToLower(strings.TrimSpace(arr[0]))
		if str == "pass" {
			status = true
		} else {
			status = false
		}

		caseStr := commonUtils.Base(strings.TrimSpace(arr[1]))
		name := strings.Replace(caseStr, path.Ext(caseStr), "", -1)
		caseIdStr := strings.Split(name, "-")[1]
		caseId, _ := strconv.Atoi(caseIdStr)

		cases[caseId] = status
	}
	//requestObj["cases"] = cases

	reqStr, _ := json.Marshal(requestObj)
	print2.PrintToCmd(string(reqStr))

	url := config.Url
	url = commonUtils.UpdateUrl(url)
	_, _ = httpClient.Post(url+constant.UrlSubmitResult, requestObj)

	//if err == nil {
	//	if pass {
	//		utils.PrintToCmd("success to submit the results")
	//	}
	//}
}
