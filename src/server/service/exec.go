package service

import (
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/server/domain"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
)

type ExecService struct {
}

func NewExecService() *ExecService {
	return &ExecService{}
}

func (s *ExecService) Exec(build domain.Build) (reply domain.OptResult) {
	serverVerbose := vari.Verbose
	vari.Verbose = build.Debug

	vari.ServerWorkDir = build.WorkDir
	if vari.ServerWorkDir != "" {
		vari.ServerWorkDir = fileUtils.AddPathSepIfNeeded(vari.ServerWorkDir)
	}

	if stringUtils.FindInArr(build.UnitTestType, constant.UnitTestTypes) { // unit test
		vari.ProductId = build.ProductId

		vari.UnitTestType = build.UnitTestType
		vari.UnitTestTool = build.UnitTestTool

		action.RunUnitTest(build.UnitTestCmd)

	} else { // ztf functional test
		vari.ProductId = build.ProductId

		action.RunZTFTest(build.Files, build.SuiteId, build.TaskId)
	}

	vari.Verbose = serverVerbose

	return
}
