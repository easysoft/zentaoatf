package service

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type TestFilterService struct {
	WorkspaceRepo *repo.WorkspaceRepo `inject:""`
	SiteService   *SiteService        `inject:""`
}

func NewTestFilterService() *TestFilterService {
	return &TestFilterService{}
}

func (s *TestFilterService) ListFilterItems(filerType commConsts.ScriptFilterType,
	siteId int, productId int) (ret []serverDomain.FilterItem, err error) {

	if filerType == commConsts.FilterWorkspace {
		ret, err = s.ListWorkspaceFilter(siteId, productId)
		return
	}

	//site, _ := s.SiteService.GetDomainObject(uint(siteId))
	//config := commDomain.WorkspaceConf{
	//	Url:      site.Url,
	//	Username: site.Username,
	//	Password: site.Password,
	//}

	if filerType == commConsts.FilterModule {

	} else if filerType == commConsts.FilterSuite {

	} else if filerType == commConsts.FilterTask {

	}

	return
}

func (s *TestFilterService) ListWorkspaceFilter(siteId int, productId int) (ret []serverDomain.FilterItem, err error) {
	workspaces, err := s.WorkspaceRepo.ListWorkspacesByProduct(siteId, productId)

	for _, item := range workspaces {
		filterItem := serverDomain.FilterItem{Label: item.Name, Value: int(item.ID)}
		ret = append(ret, filterItem)
	}

	return
}
