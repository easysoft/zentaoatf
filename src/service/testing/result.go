package testingService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	httpClient "github.com/easysoft/zentaoatf/src/service/client"
	scriptService "github.com/easysoft/zentaoatf/src/service/script"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	printUtils "github.com/easysoft/zentaoatf/src/utils/print"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"os"
	"strconv"
	"strings"
)

func SubmitResult(assert string, date string) {
	conf := configUtils.ReadCurrConfig()

	report := GetTestTestReportForSubmit(assert, date)

	for _, cs := range report.Cases {
		id := cs.Id
		runId := cs.IdInTask

		var uri string
		uri = fmt.Sprintf("testtask-runCase-%d-%d-1.json", runId, id)

		requestObj := map[string]string{"case": strconv.Itoa(id), "version": "0"}

		for _, step := range cs.Steps {
			var stepStatus string
			if step.Status {
				stepStatus = constant.PASS.String()
			} else {
				stepStatus = constant.FAIL.String()
			}

			stepResults := ""
			for _, checkpoint := range step.CheckPoints {
				stepResults += checkpoint.Actual // strconv.FormatBool(checkpoint.Status) + ": " + checkpoint.Actual
			}

			requestObj["steps["+strconv.Itoa(step.Id)+"]"] = stepStatus
			requestObj["reals["+strconv.Itoa(step.Id)+"]"] = stepResults
		}

		reqStr, _ := json.Marshal(requestObj)
		printUtils.PrintToCmd(string(reqStr))

		url := conf.Url + uri
		_, ok := httpClient.PostStr(url, requestObj)
		if ok {
			printUtils.PrintToCmd(fmt.Sprintf("success to submit the results for case %d", id))
		}
	}
}

func GetTestTestReportForSubmit(assert string, date string) model.TestReport {
	mode, name := scriptService.GetRunModeAndName(assert)
	resultPath := vari.Prefer.WorkDir + constant.LogDir + scriptService.LogFolder(mode, name, date) +
		string(os.PathSeparator) + "result.json"

	content := fileUtils.ReadFile(resultPath)
	content = strings.Replace(content, "\n", "", -1)

	var report model.TestReport
	json.Unmarshal([]byte(content), &report)

	return report
}
