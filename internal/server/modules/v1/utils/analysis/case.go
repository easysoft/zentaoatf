package analysisUtils

import (
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"path"
	"strings"
)

func FilterCaseByResult(cases []string, req serverDomain.WsReq) []string { // scope: all | fail
	scope := req.Scope

	report, err := GetReport(req.ProjectPath, req.Seq)

	cases := make([]string, 0)

	for _, cs := range report.FuncResult {
		if cs.Status != "pass" {
			cases = append(cases, cs.Path)
		}
	}

	return cases
}
