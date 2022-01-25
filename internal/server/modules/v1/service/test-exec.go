package service

import (
	"encoding/json"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/exec"
	"github.com/jinzhu/copier"
	"path/filepath"
)

type TestExecService struct {
	TestExecRepo *repo.TestExecRepo `inject:""`
}

func NewTestExecService() *TestExecService {
	return &TestExecService{}
}

func (s *TestExecService) List(projectPath string) (ret []serverDomain.TestReportSummary, err error) {
	reportFiles := scriptUtils.ListReport(projectPath)

	dir := filepath.Join(projectPath, commConsts.LogDirName)

	for _, seq := range reportFiles {
		pth := filepath.Join(dir, seq, commConsts.ResultJson)

		content := fileUtils.ReadFileBuf(pth)
		var report commDomain.ZtfReport
		err1 := json.Unmarshal(content, &report)
		if err1 != nil {
			continue
		}

		var summary serverDomain.TestReportSummary
		copier.Copy(&summary, report)
		summary.Seq = seq
		ret = append(ret, summary)
	}

	return
}

func (s *TestExecService) Get(projectPath string, seq string) (report commDomain.ZtfReport, err error) {
	dir := filepath.Join(projectPath, commConsts.LogDirName)
	pth := filepath.Join(dir, seq, commConsts.ResultJson)

	content := fileUtils.ReadFileBuf(pth)
	err = json.Unmarshal(content, &report)

	return
}

func (s *TestExecService) Delete(projectPath string, seq string) (err error) {
	dir := filepath.Join(projectPath, commConsts.LogDirName)

	di := filepath.Join(dir, seq)
	err = fileUtils.RmDir(di)

	return
}
