package zentaoService

import (
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
)

func CommitZtfTestResult(resultDir string, noNeedConfirm bool) {
	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	report := testingService.GetZtfTestReportForSubmit(resultDir)
	CommitTestResult(report)
}

//func GetLastResult(baseUrl string, caseId int) int {
//	params := ""
//	if vari.RequestType == constant.RequestTypePathInfo {
//		params = fmt.Sprintf("0-%d-1.json", caseId)
//	} else {
//		params = fmt.Sprintf("caseID=%d", caseId)
//	}
//
//	url := baseUrl + zentaoUtils.GenApiUri("testtask", "results", params)
//	dataStr, ok := client.Get(url)
//
//	resultId := -1
//	if ok {
//		jsonData, err := simplejson.NewJson([]byte(dataStr))
//		if err == nil {
//			results, _ := jsonData.Get("results").Map()
//
//			for key, _ := range results {
//				numb, _ := strconv.Atoi(key)
//
//				if resultId < numb {
//					resultId = numb
//				}
//			}
//		}
//	}
//
//	return resultId
//}
