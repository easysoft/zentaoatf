package service

import (
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/server/domain"
	serverUtils "github.com/easysoft/zentaoatf/src/server/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"strings"
)

type ExecService struct {
}

func NewExecService() *ExecService {
	return &ExecService{}
}

func (s *ExecService) Exec(build domain.Build) (reply domain.OptResult) {
	serverVerbose := vari.Verbose
	vari.Verbose = build.Debug
	vari.RunMode = constant.RunModeRequest
	defer rollback(serverVerbose)

	s.prepareCodes(&build)
	s.prepareDir(&build)

	resultDir := ""
	if stringUtils.FindInArr(build.UnitTestType, constant.UnitTestTypes) { // unit test
		vari.ProductId = build.ProductId

		vari.UnitTestType = build.UnitTestType
		vari.UnitTestTool = build.UnitTestTool

		resultDir = action.RunUnitTest(build.UnitTestCmd)

	} else { // ztf functional test
		vari.ProductId = build.ProductId

		action.RunZTFTest(build.Files, build.SuiteId, build.TaskId)
		resultDir = vari.LogDir
	}

	serverUtils.BakLog(resultDir)
	return
}

func (s *ExecService) prepareCodes(build *domain.Build) {
	if build.WorkDir != "" {
		build.WorkDir = fileUtils.AddPathSepIfNeeded(build.WorkDir)
	}

	if build.ScmAddress != "" { // git
		serverUtils.CheckoutCodes(build)

	} else if strings.Index(build.ScriptUrl, "http") == 0 { // zip
		serverUtils.DownloadCodes(build)

	} else { // folder
		if build.ScriptUrl != "" {
			build.ScriptUrl = fileUtils.AddPathSepIfNeeded(build.ScriptUrl)
		}
		build.ProjectDir = build.ScriptUrl
	}
}

func (s *ExecService) prepareDir(build *domain.Build) {
	vari.ServerWorkDir = build.WorkDir
	vari.ServerProjectDir = build.ProjectDir

	if vari.ServerProjectDir == "" && vari.ServerWorkDir != "" {
		vari.ServerProjectDir = vari.ServerWorkDir
	} else if vari.ServerProjectDir != "" && vari.ServerWorkDir == "" {
		vari.ServerWorkDir = vari.ServerProjectDir
	} else if vari.ServerProjectDir == "" && vari.ServerWorkDir == "" {
		vari.ServerWorkDir = fileUtils.AbsolutePath(".")
		vari.ServerProjectDir = vari.ServerWorkDir
	}

	if vari.ServerWorkDir != "" {
		vari.ServerWorkDir = fileUtils.AddPathSepIfNeeded(vari.ServerWorkDir)
	}
	if vari.ServerProjectDir != "" {
		vari.ServerProjectDir = fileUtils.AddPathSepIfNeeded(vari.ServerProjectDir)
	}
}

func rollback(serverVerbose bool) {
	vari.Verbose = serverVerbose
	vari.RunMode = constant.RunModeCommon
}
