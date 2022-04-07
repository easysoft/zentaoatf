package service

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	codeHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/code"
	"github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	zentaoHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"github.com/kataras/iris/v12"
)

type TestScriptService struct {
	WorkspaceRepo     *repo.WorkspaceRepo `inject:""`
	SiteService       *SiteService        `inject:""`
	TestResultService *TestResultService  `inject:""`
}

func NewTestScriptService() *TestScriptService {
	return &TestScriptService{}
}

func (s *TestScriptService) LoadTestScriptsBySiteProduct(
	siteId, productId int, filerType string, filerValue int) (root serverDomain.TestAsset, err error) {

	scriptIdsFromZentao := s.getScriptIdsFromZentao(siteId, productId, filerType, filerValue)
	workspaces, _ := s.WorkspaceRepo.ListByProduct(siteId, productId)

	// load scripts from disk
	root = serverDomain.TestAsset{Path: "", Title: "测试脚本", Type: commConsts.Root, Slots: iris.Map{"icon": "icon"}}
	for _, workspace := range workspaces {

		if filerType == string(commConsts.FilterWorkspace) &&
			(filerValue > 0 && uint(filerValue) != workspace.ID) { // filter by workspace
			continue
		}

		var scriptsInDir serverDomain.TestAsset

		if workspace.Type == commConsts.ZTF {
			scriptsInDir, _ = scriptHelper.LoadScriptTree(workspace, scriptIdsFromZentao)
		} else {
			scriptsInDir, _ = codeHelper.LoadCodeTree(workspace)
		}

		root.Children = append(root.Children, &scriptsInDir)
	}

	return
}

func (s *TestScriptService) LoadCodeChildren(dir string, workspaceId int) (nodes []*serverDomain.TestAsset, err error) {
	workspace, _ := s.WorkspaceRepo.Get(uint(workspaceId))

	nodes, _ = codeHelper.LoadCodeNodesInDir(dir, workspaceId, workspace.Type)

	return
}

func (s *TestScriptService) getScriptIdsFromZentao(siteId, productId int, filerType string, filerValue int) (
	ret map[int]string) {

	if filerType == "" || filerValue < 0 {
		return nil
	}

	// get script ids from zentao
	currSite, _ := s.SiteService.Get(uint(siteId))
	config := commDomain.WorkspaceConf{
		Url:      currSite.Url,
		Username: currSite.Username,
		Password: currSite.Password,
	}

	if filerType == string(commConsts.FilterModule) {
		ret, _ = zentaoHelper.GetCaseIdsInZentaoModule(productId, filerValue, config)
	} else if filerType == string(commConsts.FilterSuite) {
		ret, _ = zentaoHelper.GetCaseIdsInZentaoSuite(productId, filerValue, config)
	} else if filerType == string(commConsts.FilterTask) {
		ret, _ = zentaoHelper.GetCaseIdsInZentaoTask(productId, filerValue, config)
	}

	return
}

func (s *TestScriptService) GetCaseIdsFromReport(workspaceId int, seq, scope string) (caseIds []string, err error) {
	report, err := s.TestResultService.Get(workspaceId, seq)

	if err != nil {
		return
	}

	for _, cs := range report.FuncResult {
		path := cs.Path
		status := cs.Status

		if path != "" && (scope == "all" || scope == status.String()) {
			caseIds = append(caseIds, path)
		}
	}

	return
}

func (s *TestScriptService) UpdateCode(script serverDomain.TestScript) (err error) {
	fileUtils.WriteFile(script.Path, script.Code)

	return
}
