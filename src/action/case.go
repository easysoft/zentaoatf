package action

import (
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	assertUtils "github.com/easysoft/zentaoatf/src/utils/assert"
	scriptUtils "github.com/easysoft/zentaoatf/src/utils/script"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
)

func CommitCases(files []string) {
	cases := assertUtils.GetCaseByDirAndFile(files)

	for _, cs := range cases {
		pass, id, _, title := zentaoUtils.GetCaseInfo(cs)

		if pass {
			stepMap, stepTypeMap, expectMap := scriptUtils.GetStepAndExpectMap(cs)

			isIndependent, expectIndependentContent := zentaoUtils.GetDependentExpect(cs)
			if isIndependent {
				expectMap = scriptUtils.GetExpectMapFromIndependentFile(expectMap, expectIndependentContent, true)
			}

			zentaoService.CommitCase(id, title, stepMap, stepTypeMap, expectMap)
		}
	}
}
