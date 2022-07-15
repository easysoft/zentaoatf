package analysisHelper

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
)

func FilterCaseByResult(cases []string, req serverDomain.TestSet) (ret []string) { // scope: all | fail
	scope := req.Scope
	caseIdMap, _ := getCaseIdMapFromReport(req)

	for _, cs := range cases {
		if scope.String() == "all" || (scope == caseIdMap[cs]) {
			ret = append(ret, cs)
		}
	}

	return
}

func getCaseIdMapFromReport(req serverDomain.TestSet) (ret map[string]commConsts.ResultStatus, err error) {
	report, _, err := ReadReportByWorkspaceSeq(req.WorkspacePath, req.Seq)
	if err != nil {
		logUtils.Errorf("fail to get case ids for %s %s", req.WorkspacePath, req.Seq)
		return
	}

	ret = map[string]commConsts.ResultStatus{}
	for _, cs := range report.FuncResult {
		ret[cs.Path] = cs.Status
	}

	return
}
