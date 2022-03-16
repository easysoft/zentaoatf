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

func (s *SiteService) Create(site model.Site) (id uint, err error) {
	id, err = s.SiteRepo.Create(site)

	//err = configUtils.UpdateSite(site, workspacePath)

	return
}

func (s *SiteService) Update(site model.Site) (err error) {
	err = s.SiteRepo.Update(site)

	//err = configUtils.UpdateSite(site, workspacePath)
	return
}

func (s *SiteService) Delete(id uint) error {
	return s.SiteRepo.Delete(id)
}

func (s *SiteService) LoadSites(currSiteId int) (sites []serverDomain.ZentaoSite, currSite serverDomain.ZentaoSite, err error) {
	req := serverDomain.ReqPaginate{PaginateReq: domain.PaginateReq{Page: 1, PageSize: 10000}}
	pageData, err := s.Paginate(req)
	if err != nil {
		return
	}

	pos := pageData.Result.([]*model.Site)

	var first serverDomain.ZentaoSite
	for idx, item := range pos {
		site := serverDomain.ZentaoSite{
			Id:       int(item.ID),
			Url:      item.Url,
			Username: item.Username,
			Password: item.Password,
		}

		if uint(currSiteId) == item.ID {
			currSite = site
		}

		if idx == 0 {
			first = site
		}

		sites = append(sites, site)
	}

	if currSite.Id == 0 { // not found, use the first one
		currSite = first
	}

	return
}
