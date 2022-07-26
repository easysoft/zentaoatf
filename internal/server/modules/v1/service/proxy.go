package service

import (
	"errors"

	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
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

func (s *ProxyService) Get(id uint) (proxy model.Proxy, err error) {
	return s.ProxyRepo.Get(id)
}

func (s *ProxyService) Create(proxy model.Proxy) (id uint, err error) {
	proxy.Path = zentaoHelper.FixSiteUrl(proxy.Path)
	if proxy.Path == "" {
		err = errors.New("url not right")
		return
	}
	proxy.Path = fileUtils.AddUrlPathSepIfNeeded(proxy.Path)
	err = s.CheckServer(proxy.Path)
	if err != nil {
		return
	}
	id, err = s.ProxyRepo.Create(proxy)
	return
}

func (s *ProxyService) Update(proxy model.Proxy) (err error) {
	proxy.Path = zentaoHelper.FixSiteUrl(proxy.Path)
	if proxy.Path == "" {
		err = errors.New("url not right")
		return
	}
	proxy.Path = fileUtils.AddUrlPathSepIfNeeded(proxy.Path)
	err = s.CheckServer(proxy.Path)
	if err != nil {
		return err
	}
	err = s.ProxyRepo.Update(proxy)
	return
}

func (s *ProxyService) Delete(id uint) error {
	return s.ProxyRepo.Delete(id)
}

func (s *ProxyService) CheckServer(url string) (err error) {
	_, err = httpUtils.Get(url + "api/v1/heartbeat")
	return
}
