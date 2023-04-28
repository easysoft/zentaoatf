package analysisHelper

import (
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
)

func GetMethodByLineNum(num int, req map[]) (ret []string) { // scope: all | fail
	scope := req.Scope
	caseIdMap, _ := getCaseIdMapFromReport(req)

	for _, cs := range cases {
		if scope.String() == "all" || (scope == caseIdMap[cs]) {
			ret = append(ret, cs)
		}
	}

	return
}
