package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type SiteService struct {
	SiteRepo *repo.SiteRepo `inject:""`
}

func NewSiteService() *SiteService {
	return &SiteService{}
}

func (s *SiteService) Paginate(req serverDomain.ReqPaginate) (ret domain.PageData, err error) {
	ret, err = s.SiteRepo.Paginate(req)
	return
}

func (s *SiteService) Get(id uint) (site model.Site, err error) {
	return s.SiteRepo.Get(id)
}

func (s *SiteService) Create(site model.Site) (uint, error) {
	return s.SiteRepo.Create(site)
}

func (s *SiteService) Update(site model.Site) error {
	return s.SiteRepo.Update(site)
}

func (s *SiteService) Delete(id uint) error {
	return s.SiteRepo.Delete(id)
}
