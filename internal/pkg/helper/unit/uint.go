package unitHelper

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
)

func GetUnitTools(args []string, startIndex int) {
	if startIndex > len(args)-1 { // no exec command
		return
	}

	str := args[startIndex]

	if str == commConsts.Maven.String() {
		commConsts.UnitBuildTool = commConsts.Maven
	} else if str == commConsts.Mocha.String() {
		commConsts.UnitBuildTool = commConsts.Mocha
	} else if str == commConsts.RobotFramework.String() {
		commConsts.UnitTestTool = commConsts.RobotFramework
	}

	if commConsts.UnitTestTool == "" {
		commConsts.UnitTestTool = commConsts.TestTool(commConsts.UnitTestType)
	}
}
