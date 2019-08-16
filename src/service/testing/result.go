package testingService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	scriptService "github.com/easysoft/zentaoatf/src/service/script"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"os"
	"strings"
)

func GetTestTestReportForSubmit(scriptFile string, resultDate string) model.TestReport {
	mode, name := scriptService.GetRunModeAndName(scriptFile)
	resultPath := vari.Prefer.WorkDir + constant.LogDir + scriptService.LogFolder(mode, name, resultDate) +
		string(os.PathSeparator) + "result.json"

	content := fileUtils.ReadFile(resultPath)
	content = strings.Replace(content, "\n", "", -1)

	var report model.TestReport
	json.Unmarshal([]byte(content), &report)

	return report
}

func SaveTestTestReportAfterSubmit(scriptFile string, resultDate string, content string) {
	mode, name := scriptService.GetRunModeAndName(scriptFile)
	resultPath := vari.Prefer.WorkDir + constant.LogDir + scriptService.LogFolder(mode, name, resultDate) +
		string(os.PathSeparator) + "result.json"

	fileUtils.WriteFile(resultPath, content)
}

func GetStepContent(step model.StepLog) string {
	var stepsContent string
	if !vari.RunFromCui {
		stepsContent = GetStepHtml(step)
	} else {
		stepsContent = GetStepText(step)
	}

	return stepsContent
}

func GetStepHtml(step model.StepLog) string {
	stepResults := make([]string, 0)

	stepStatus := stringUtils.BoolToPass(step.Status)

	stepTxt := fmt.Sprintf(
		"<p><b>%s: %s</b></p>",
		step.Name, stepStatus)

	for _, checkpoint := range step.CheckPoints {
		checkpointStatus := stringUtils.BoolToPass(checkpoint.Status)

		text := fmt.Sprintf(
			"<p>&nbsp;Checkpoint: %s</p>"+
				"<p>&nbsp;&nbsp;Expect</p>"+
				"&nbsp;&nbsp;&nbsp;%s"+
				"<p>&nbsp;&nbsp;Actual<p/>"+
				"&nbsp;&nbsp;&nbsp;%s",
			checkpointStatus, checkpoint.Expect, checkpoint.Actual)

		stepResults = append(stepResults, text)
	}

	return stepTxt + strings.Join(stepResults, "<br/>")
}

func GetStepText(step model.StepLog) string {
	stepResults := make([]string, 0)

	stepStatus := stringUtils.BoolToPass(step.Status)

	stepTxt := fmt.Sprintf(
		"%s: %s\n",
		step.Name, stepStatus)

	for _, checkpoint := range step.CheckPoints {
		checkpointStatus := stringUtils.BoolToPass(checkpoint.Status)

		text := fmt.Sprintf(
			" Checkpoint: %s\n"+
				"  Expect\n"+
				"   %s\n"+
				"  Actual\n"+
				"   %s",
			checkpointStatus, checkpoint.Expect, checkpoint.Actual)

		stepResults = append(stepResults, text)
	}

	return stepTxt + strings.Join(stepResults, "\n")
}
