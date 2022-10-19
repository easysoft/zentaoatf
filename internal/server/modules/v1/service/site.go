package service

import (
	"errors"
	"strings"

	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	"github.com/easysoft/zentaoatf/pkg/domain"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
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

func (s *SiteService) Create(site model.Site) (id uint, isDuplicate bool, err error) {
	site.Url = zentaoHelper.FixSiteUrl(site.Url)
	if site.Url == "" {
		err = errors.New("url not right")
		return
	}

	site.Url = fileUtils.AddUrlPathSepIfNeeded(site.Url)

	config := configHelper.LoadBySite(site)
	err = zentaoHelper.Login(config)
	if err != nil {
		config.Url += "zentao/"
		site.Url += "zentao/"
		err = zentaoHelper.Login(config)
		if err != nil {
			return
		}
	}

	id, isDuplicate, err = s.SiteRepo.Create(&site)

	return
}

func (s *SiteService) Update(site model.Site) (isDuplicate bool, err error) {
	site.Url = zentaoHelper.FixSiteUrl(site.Url)
	if site.Url == "" {
		err = errors.New("url not right")
		return
	}

	site.Url = fileUtils.AddUrlPathSepIfNeeded(site.Url)

	config := configHelper.LoadBySite(site)
	err = zentaoHelper.Login(config)
	if err != nil {
		config.Url += "zentao/"
		site.Url += "zentao/"
		err = zentaoHelper.Login(config)
		if err != nil {
			return
		}
	}

	isDuplicate, err = s.SiteRepo.Update(site)
	if isDuplicate || err != nil {
		return
	}

	workspaces, _ := s.WorkspaceRepo.ListBySite(site.ID)
	for _, item := range workspaces {
		s.WorkspaceService.UpdateConfig(item, "site")
	}

	return
}

func (s *SiteService) Delete(id uint) error {
	err := s.WorkspaceRepo.DeleteBySite(id)
	if err != nil {
		return err
	}

	return s.SiteRepo.Delete(id)
}

func (s *SiteService) LoadSites(currSiteId int, lang string) (sites []serverDomain.ZentaoSite, currSite serverDomain.ZentaoSite, err error) {
	req := serverDomain.ReqPaginate{PaginateReq: domain.PaginateReq{Page: 1, PageSize: 10000}}
	pageData, err := s.Paginate(req)
	if err != nil {
		return
	}

	pos := pageData.Result.([]*model.Site)
	if len(pos) == 0 {
		s.CreateEmptySite(lang)
		pageData, err = s.Paginate(req)
		pos = pageData.Result.([]*model.Site)
	}

	sites = []serverDomain.ZentaoSite{}
	currIndex := 0
	for idx, item := range pos {
		site := serverDomain.ZentaoSite{
			Id:       int(item.ID),
			Name:     item.Name,
			Url:      item.Url,
			Username: item.Username,
			Password: item.Password,
		}

		if uint(currSiteId) == item.ID {
			currIndex = idx
		}

		sites = append(sites, site)
	}

	currSite = sites[currIndex] // default is first one

	return
}

func (s *SiteService) CreateEmptySite(lang string) (err error) {
	name := "Local"
	if strings.Index(strings.ToLower(lang), "zh") > -1 {
		name = "本地"
	}

	po := model.Site{
		Name: name,
		Url:  "",
	}
	_, _, err = s.SiteRepo.Create(&po)

	return
}
