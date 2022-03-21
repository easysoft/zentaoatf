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

	workspaces, _ := s.WorkspaceRepo.ListWorkspacesByProduct(siteId, productId)

	// load scripts from disk
	root = serverDomain.TestAsset{Path: "", Title: "测试脚本", Type: commConsts.Root, Slots: iris.Map{"icon": "icon"}}
	for _, workspace := range workspaces {
		if workspace.Type == commConsts.ZTF {
			if filerType == string(commConsts.FilterWorkspace) && uint(filerValue) != workspace.ID { // filter by workspace
				continue
			}

			scriptsInDir, _ := scriptUtils.LoadScriptTree(workspace.Path)

			root.Children = append(root.Children, &scriptsInDir)
		}
	}

	if filerType == string(commConsts.FilterWorkspace) {
		return
	}

	// get script ids from zentao
	currSite, _ := s.SiteService.Get(uint(siteId))
	config := commDomain.WorkspaceConf{
		Url:      currSite.Url,
		Username: currSite.Username,
		Password: currSite.Password,
	}

	caseIdsInZentao := map[int]string{}
	if filerType == string(commConsts.FilterModule) {
		caseIdsInZentao, _ = zentaoHelper.GetCaseIdsInZentaoModule(productId, filerValue, config)
	} else if filerType == string(commConsts.FilterWorkspace) {

	} else if filerType == string(commConsts.FilterWorkspace) {

	}

	// filter scripts by the zentao ids
	root = s.filterScriptsByCaseId(root, caseIdsInZentao)

	return
}

func (s *TestScriptService) filterScriptsByCaseId(root serverDomain.TestAsset, caseIdsInZentao map[int]string) (
	ret serverDomain.TestAsset) {

	ret = root

	return
}
