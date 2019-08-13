package biz

import (
	"encoding/json"
	httpClient "github.com/easysoft/zentaoatf/src/client"
	"github.com/easysoft/zentaoatf/src/utils"
	"path"
	"strconv"
	"strings"
)

func SubmitResult(caseList []string) {
	config := utils.ReadCurrConfig()

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

		caseStr := utils.Base(strings.TrimSpace(arr[1]))
		name := strings.Replace(caseStr, path.Ext(caseStr), "", -1)
		caseIdStr := strings.Split(name, "-")[1]
		caseId, _ := strconv.Atoi(caseIdStr)

		cases[caseId] = status
	}
	//requestObj["cases"] = cases

	reqStr, _ := json.Marshal(requestObj)
	utils.PrintToCmd(string(reqStr))

	url := config.Url
	url = utils.UpdateUrl(url)
	_, _ = httpClient.Post(url+utils.UrlSubmitResult, requestObj)

	//if err == nil {
	//	if pass {
	//		utils.PrintToCmd("success to submit the results")
	//	}
	//}
}
