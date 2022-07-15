package service

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
)

type TestFilterService struct {
	WorkspaceRepo *repo.WorkspaceRepo `inject:""`
	SiteService   *SiteService        `inject:""`
}

func NewTestFilterService() *TestFilterService {
	return &TestFilterService{}
}

func (s *TestFilterService) ListFilterItems(filerType commConsts.ScriptFilterType,
	siteId, productId uint) (ret interface{}, err error) {

	if filerType == commConsts.FilterWorkspace {
		ret, err = s.ListWorkspaceFilter(siteId, productId)
		return
	}

	site, _ := s.SiteService.GetDomainObject(siteId)
	config := commDomain.WorkspaceConf{
		Url:      site.Url,
		Username: site.Username,
		Password: site.Password,
	}

	if filerType == commConsts.FilterModule {
		ret, err = s.ListModuleFilter(config, productId)
	} else if filerType == commConsts.FilterSuite {
		ret, err = s.ListSuiteFilter(config, productId)
	} else if filerType == commConsts.FilterTask {
		ret, err = s.ListTaskFilter(config, productId)
	}

	return
}

func (s *TestFilterService) ListWorkspaceFilter(siteId, productId uint) (ret []serverDomain.FilterItem, err error) {
	workspaces, err := s.WorkspaceRepo.ListByProduct(siteId, productId)

	for _, item := range workspaces {
		filterItem := serverDomain.FilterItem{Label: item.Name, Value: int(item.ID)}
		ret = append(ret, filterItem)
	}

	return
}

func (s *TestFilterService) ListModuleFilter(config commDomain.WorkspaceConf, productId uint) (ret []serverDomain.FilterItem, err error) {
	modules, _ := zentaoHelper.LoadCaseModule(productId, config)

	for _, item := range modules {
		filterItem := serverDomain.FilterItem{Label: item.Name, Value: item.Id}
		ret = append(ret, filterItem)
	}

	return
}

func (s *TestFilterService) ListSuiteFilter(config commDomain.WorkspaceConf, productId uint) (ret []serverDomain.FilterItem, err error) {
	suites, _ := zentaoHelper.LoadSuite(productId, config)

	for _, item := range suites {
		filterItem := serverDomain.FilterItem{Label: item.Name, Value: item.Id}
		ret = append(ret, filterItem)
	}

	return
}
func (s *TestFilterService) ListTaskFilter(config commDomain.WorkspaceConf, productId uint) (ret []serverDomain.FilterItem, err error) {
	tasks, _ := zentaoHelper.LoadTask(productId, config)

	for _, item := range tasks {
		filterItem := serverDomain.FilterItem{Label: item.Name, Value: item.Id}
		ret = append(ret, filterItem)
	}

	return
}
