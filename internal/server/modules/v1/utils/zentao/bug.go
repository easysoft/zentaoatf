package zentaoUtils

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/domain"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/analysis"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/config"
	"github.com/fatih/color"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
)

func CommitBug(bug commDomain.ZtfBug, projectPath string) (err error) {
	config := configUtils.LoadByProjectPath(projectPath)
	Login(config)

	bug.Steps = strings.Replace(bug.Steps, " ", "&nbsp;", -1)

	// bug-create-1-0-caseID=1,version=3,resultID=93,runID=0,stepIdList=9_12_
	// bug-create-1-0-caseID=1,version=3,resultID=84,runID=6,stepIdList=9_12_,testtask=2,projectID=1,buildID=1
	extras := fmt.Sprintf("caseID=%s,version=%s,stepIdList=%s",
		bug.Case, bug.Version, bug.StepIds)

	// $productID, $branch = '', $extras = ''
	params := ""
	if commConsts.RequestType == commConsts.PathInfo {
		params = fmt.Sprintf("%s-0-%s", bug.Product, extras)
	} else {
		params = fmt.Sprintf("productID=%s&branch=0&$extras=%s", bug.Product, extras)
	}
	//params = ""
	url := config.Url + GenApiUri("bug", "create", params)

	_, ok := httpUtils.Post(url, bug, false)

	msg := ""
	if ok {
		msg = i118Utils.Sprintf("success_to_report_bug", bug.Case)
	} else {
		msg = color.RedString(err.Error())
	}
	logUtils.Info(msg)

	return
}

func PrepareBug(projectPath, seq string, caseIdStr string) (bug commDomain.ZtfBug) {
	caseId, err := strconv.Atoi(caseIdStr)

	if err != nil {
		return
	}

	report, err := analysisUtils.ReadReport(projectPath, seq)
	if err != nil {
		return
	}

	for _, cs := range report.FuncResult {
		if cs.Id != caseId {
			continue
		}

		steps := make([]string, 0)
		stepIds := ""
		for _, step := range cs.Steps {
			if step.Status == commConsts.FAIL {
				stepIds += step.Id + "_"
			}

			stepsContent := GetStepText(step)
			steps = append(steps, stepsContent)
		}

		bug = commDomain.ZtfBug{
			Title:   cs.Title,
			Product: strconv.Itoa(cs.ProductId), Case: strconv.Itoa(cs.Id),
			Uid: uuid.NewV4().String(), CaseVersion: "0", OldTaskID: "0",
			Steps: strings.Join(steps, "\n"), StepIds: stepIds,
			Version: "trunk", Severity: "3", Pri: "3", OpenedBuild: map[string]string{"0": "trunk"},
		}
		return
	}

	return
}
