package unitHelper

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
)

func GetUnitTools(args []string, startIndex int) {
	if startIndex > len(args)-1 { // no exec command
		return
	}

	str := args[startIndex]

	if str == commConsts.UnitTestToolMvn {
		commConsts.UnitBuildTool = commConsts.Maven
	} else if str == commConsts.UnitTestToolMocha {
		commConsts.UnitBuildTool = commConsts.Mocha
	} else if str == commConsts.UnitTestToolRobot {
		commConsts.UnitTestTool = commConsts.RobotFramework
	}

	if commConsts.UnitTestTool == "" {
		commConsts.UnitTestTool = commConsts.TestTool(commConsts.UnitTestType)
	}
}
