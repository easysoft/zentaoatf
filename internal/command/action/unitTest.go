package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
)

func RunUnitTest(cmdStr string, taskIdOrName string) {
	testSet := serverDomain.TestSet{
		ProductId:     stringUtils.ParseInt(commConsts.ProductId),
		WorkspacePath: commConsts.WorkDir,

		Cmd:       cmdStr,
		TestTool:  commConsts.UnitTestTool,
		BuildTool: commConsts.UnitBuildTool,
	}
	if testSet.ProductId != 0 {
		testSet.SubmitResult = true
	}

	if taskIdOrName != "" {
		testSet.Name = taskIdOrName
	}

	req := serverDomain.ExecReq{
		Act:      commConsts.ExecUnit,
		TestSets: []serverDomain.TestSet{testSet},
	}

	execHelper.Exec(nil, req, nil)
}
