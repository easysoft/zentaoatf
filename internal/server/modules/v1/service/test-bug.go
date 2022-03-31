package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	zentaoHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type TestBugService struct {
	SiteRepo *repo.SiteRepo `inject:""`
}

func NewTestBugService() *TestBugService {
	return &TestBugService{}
}

func (s *TestBugService) Submit(bug commDomain.ZtfBug, siteId, productId int) (err error) {
	site, err := s.SiteRepo.Get(uint(siteId))
	config := configUtils.LoadBySite(site)

	bug.Product = productId
	err = zentaoHelper.CommitBug(bug, config)

	return
}

func (s *TestBugService) GetBugFields(siteId, productId int) (bugFields commDomain.ZentaoBugFields, err error) {
	site, err := s.SiteRepo.Get(uint(siteId))
	config := configUtils.LoadBySite(site)

	bugFields, err = zentaoHelper.GetBugFiledOptions(config, productId)
	return
}
