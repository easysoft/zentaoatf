package unitHelper

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"strings"
)

func GetUnitTools(args []string, startIndex int) {
	str := args[startIndex]

	if commConsts.UnitTestType == commConsts.UnitTestTypeAllure {
		commConsts.UnitTestTool = commConsts.Allure
	} else if str == commConsts.UnitTestToolMvn {
		commConsts.UnitBuildTool = commConsts.Maven
	} else if str == commConsts.UnitTestToolMocha {
		commConsts.UnitTestTool = commConsts.Puppeteer
		commConsts.UnitBuildTool = commConsts.Mocha
	} else if str == commConsts.UnitTestToolRobot {
		commConsts.UnitTestTool = commConsts.RobotFramework
	} else {
		cmdStr := strings.ToLower(strings.Join(args[startIndex:], "; "))
		if strings.Index(cmdStr, commConsts.Playwright.String()) > -1 {
			commConsts.UnitTestTool = commConsts.Playwright
		} else {
			commConsts.UnitTestTool = commConsts.TestTool(str)
		}
	}
}
