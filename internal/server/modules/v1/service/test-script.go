package service

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	zentaoHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"github.com/kataras/iris/v12"
)

type TestScriptService struct {
	WorkspaceRepo *repo.WorkspaceRepo `inject:""`
	SiteService   *SiteService        `inject:""`
}

func NewTestScriptService() *TestScriptService {
	return &TestScriptService{}
}

func (s *TestScriptService) LoadTestScriptsBySiteProduct(
	siteId, productId int, filerType string, filerValue int) (root serverDomain.TestAsset, err error) {

	scriptIdsFromZentao := s.getScriptIdsFromZentao(siteId, productId, filerType, filerValue)
	workspaces, _ := s.WorkspaceRepo.ListWorkspacesByProduct(siteId, productId)

	// load scripts from disk
	root = serverDomain.TestAsset{Path: "", Title: "测试脚本", Type: commConsts.Root, Slots: iris.Map{"icon": "icon"}}
	for _, workspace := range workspaces {
		if workspace.Type != commConsts.ZTF {
			continue
		}

		if filerType == string(commConsts.FilterWorkspace) &&
			(filerValue > 0 && uint(filerValue) != workspace.ID) { // filter by workspace
			continue
		}

		scriptsInDir, _ := scriptUtils.LoadScriptTree(workspace, scriptIdsFromZentao)

		root.Children = append(root.Children, &scriptsInDir)
	}

	if filerType == string(commConsts.FilterWorkspace) || filerValue == 0 {
		return
	}

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
