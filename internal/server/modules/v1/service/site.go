package service

import (
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/comm/helper/zentao"
	"github.com/easysoft/zentaoatf/internal/pkg/domain"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
)

type SiteService struct {
	SiteRepo         *repo.SiteRepo      `inject:""`
	WorkspaceRepo    *repo.WorkspaceRepo `inject:""`
	WorkspaceService *WorkspaceService   `inject:""`
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
func (s *SiteService) GetDomainObject(id uint) (site serverDomain.ZentaoSite, err error) {
	po, _ := s.SiteRepo.Get(id)

	site = serverDomain.ZentaoSite{
		Url:      po.Url,
		Username: po.Username,
		Password: po.Password,
	}

	return
}

func (s *SiteService) Create(site model.Site) (id uint, err error) {
	config := configHelper.LoadBySite(site)
	err = zentaoHelper.Login(config)
	if err != nil {
		return
	}

	id, err = s.SiteRepo.Create(&site)

	return
}

func (s *SiteService) Update(site model.Site) (err error) {
	config := configHelper.LoadBySite(site)
	err = zentaoHelper.Login(config)
	if err != nil {
		return
	}

	err = s.SiteRepo.Update(site)

	workspaces, _ := s.WorkspaceRepo.ListBySite(site.ID)
	for _, item := range workspaces {
		s.WorkspaceService.UpdateConfig(item, true)
	}

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
	if len(pos) == 0 {
		s.CreateEmptySite()
		pageData, err = s.Paginate(req)
		pos = pageData.Result.([]*model.Site)
	}

	sites = []serverDomain.ZentaoSite{}
	var first serverDomain.ZentaoSite
	for idx, item := range pos {
		site := serverDomain.ZentaoSite{
			Id:       int(item.ID),
			Name:     item.Name,
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

func (s *SiteService) CreateEmptySite() (err error) {
	po := model.Site{
		Name: "无站点",
		Url:  "",
	}
	_, err = s.SiteRepo.Create(&po)

	return
}
