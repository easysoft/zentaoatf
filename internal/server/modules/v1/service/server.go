package service

import (
	"errors"

	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
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
	server.Path = zentaoHelper.FixSiteUrl(server.Path)
	if server.Path == "" {
		err = errors.New(i118Utils.Sprintf("wrong_url_format"))
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
	server.Path = zentaoHelper.FixSiteUrl(server.Path)
	if server.Path == "" {
		err = errors.New(i118Utils.Sprintf("wrong_url_format"))
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
	info, _ := s.Get(id)
	if info.IsDefault {
		return errors.New(i118Utils.Sprintf("no_delete_default_server"))
	}
	return s.ServerRepo.Delete(id)
}

func (s *ServerService) CheckServer(url string) (err error) {
	_, err = httpUtils.Get(url + "api/v1/heartbeat")
	return
}
