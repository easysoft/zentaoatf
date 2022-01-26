package analysisUtils

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
)

func FilterCaseByResult(cases []string, req serverDomain.WsReq) (ret []string) { // scope: all | fail
	scope := req.Scope
	caseIdMap, _ := getCaseIdMapFromReport(req)

	for _, cs := range cases {
		if scope.String() == "all" || (scope == caseIdMap[cs]) {
			ret = append(ret, cs)
		}
	}

	return
}

func getCaseIdMapFromReport(req serverDomain.WsReq) (ret map[string]commConsts.ResultStatus, err error) {
	report, err := GetReport(req.ProjectPath, req.Seq)
	if err != nil {
		logUtils.Errorf("fail to get case ids for %s %s", req.ProjectPath, req.Seq)
		return
	}

	ret = map[string]commConsts.ResultStatus{}
	for _, cs := range report.FuncResult {
		ret[cs.Path] = cs.Status
	}

	return
}
