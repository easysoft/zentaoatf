package service

import (
	"errors"

	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
)

type ServerService struct {
	ServerRepo *repo.ServerRepo `inject:""`
}

func NewServerService() *ServerService {
	return &ServerService{}
}

func (s *ServerService) List() (ret []model.Server, err error) {
	ret, err = s.ServerRepo.List()
	return
}

func (s *ServerService) Get(id uint) (server model.Server, err error) {
	return s.ServerRepo.Get(id)
}

func (s *ServerService) Create(server model.Server) (id uint, err error) {
	server.Path = zentaoHelper.FixSiteUlt(server.Path)
	if server.Path == "" {
		err = errors.New("url not right")
		return
	}
	server.Path = fileUtils.AddUrlPathSepIfNeeded(server.Path)
	err = s.CheckServer(server.Path)
	if err != nil {
		return
	}
	id, err = s.ServerRepo.Create(server)
	return
}

func (s *ServerService) Update(server model.Server) (err error) {
	server.Path = zentaoHelper.FixSiteUlt(server.Path)
	if server.Path == "" {
		err = errors.New("url not right")
		return
	}
	server.Path = fileUtils.AddUrlPathSepIfNeeded(server.Path)
	if server.IsDefault {
		err = s.CheckServer(server.Path)
		if err != nil {
			return err
		}
	}
	err = s.ServerRepo.Update(server)
	return
}

func (s *ServerService) Delete(id uint) error {
	return s.ServerRepo.Delete(id)
}

func (s *ServerService) CheckServer(url string) (err error) {
	_, err = httpUtils.Get(url + "api/v1/heartbeat")
	return
}
