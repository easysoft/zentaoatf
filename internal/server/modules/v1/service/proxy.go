package service

import (
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
)

type ProxyService struct {
	ProxyRepo *repo.ProxyRepo `inject:""`
}

func NewProxyService() *ProxyService {
	return &ProxyService{}
}

func (s *ProxyService) List() (ret []model.Proxy, err error) {
	ret, err = s.ProxyRepo.List()
	return
}

func (s *ProxyService) Get(id uint) (site model.Proxy, err error) {
	return s.ProxyRepo.Get(id)
}

func (s *ProxyService) Create(site model.Proxy) (id uint, err error) {
	id, err = s.ProxyRepo.Create(site)
	return
}

func (s *ProxyService) Update(site model.Proxy) (err error) {

	err = s.ProxyRepo.Update(site)
	return
}

func (s *ProxyService) Delete(id uint) error {
	return s.ProxyRepo.Delete(id)
}
