package service

import (
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
)

type TestBugService struct {
	SiteRepo *repo.SiteRepo `inject:""`
}

func NewTestBugService() *TestBugService {
	return &TestBugService{}
}

func (s *TestBugService) Submit(bug commDomain.ZtfBug, siteId, productId int) (err error) {
	site, err := s.SiteRepo.Get(uint(siteId))
	config := configHelper.LoadBySite(site)

	bug.Product = productId
	err = zentaoHelper.CommitBug(bug, config)

	return
}

func (s *TestBugService) GetBugFields(siteId, productId int) (bugFields commDomain.ZentaoBugFields, err error) {
	site, err := s.SiteRepo.Get(uint(siteId))
	config := configHelper.LoadBySite(site)

	bugFields, err = zentaoHelper.GetBugFiledOptions(config, productId)
	return
}

func (s *TestBugService) LoadBugs(siteId, productId int) (bugs []commDomain.ZentaoBug, err error) {
	if siteId == 1 {
		return
	}
	site, _ := s.SiteRepo.Get(uint(siteId))
	config := configHelper.LoadBySite(site)

	allBugs, err := zentaoHelper.LoadBugs(productId, config)
	for _, bug := range allBugs {
		if bug.Case > 0 {
			bugs = append(bugs, bug)
		}
	}
	return
}
