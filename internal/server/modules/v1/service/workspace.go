package service

import (
	"errors"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type WorkspaceService struct {
	WorkspaceRepo      *repo.WorkspaceRepo `inject:""`
	SiteRepo           *repo.SiteRepo      `inject:""`
	InterpreterService *InterpreterService `inject:""`
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{}
}

func (s *WorkspaceService) Paginate(req serverDomain.WorkspaceReqPaginate) (ret domain.PageData, err error) {
	ret, err = s.WorkspaceRepo.Paginate(req)
	return
}

func (s *WorkspaceService) ListWorkspacesByProduct(siteId, productId uint) (pos []model.Workspace, err error) {
	return s.WorkspaceRepo.ListByProduct(siteId, productId)
}

func (s *WorkspaceService) Get(id uint) (model.Workspace, error) {
	return s.WorkspaceRepo.Get(id)
}
func (s *WorkspaceService) GetByPath(workspacePath string) (po model.Workspace, err error) {
	return s.WorkspaceRepo.FindByPath(workspacePath)
}

func (s *WorkspaceService) Create(workspace model.Workspace) (id uint, err error) {
	if !fileUtils.IsDir(workspace.Path) {
		err = errors.New(fmt.Sprintf("目录%s不存在", workspace.Path))
		return
	}

	id, err = s.WorkspaceRepo.Create(workspace)
	s.UpdateConfig(workspace, true)

	return
}
func (s *WorkspaceService) Update(workspace model.Workspace) (err error) {
	if !fileUtils.IsDir(workspace.Path) {
		err = errors.New(fmt.Sprintf("目录%s不存在", workspace.Path))
		return
	}

	err = s.WorkspaceRepo.Update(workspace)
	s.UpdateConfig(workspace, true)

	return
}

func (s *WorkspaceService) Delete(id uint) (err error) {
	err = s.WorkspaceRepo.Delete(id)
	if err != nil {
		return
	}

	err = s.WorkspaceRepo.SetCurrWorkspace("")

	return
}

func (s *WorkspaceService) ListByProduct(siteId, productId uint) (pos []model.Workspace, err error) {
	return s.WorkspaceRepo.ListByProduct(siteId, productId)
}

func (s *WorkspaceService) UpdateAllConfig() {
	workspaces, _ := s.WorkspaceRepo.ListWorkspace()

	for _, item := range workspaces {
		s.UpdateConfig(item, true)
	}
}

func (s *WorkspaceService) UpdateConfig(workspace model.Workspace, forceUpdate bool) (err error) {
	site, _ := s.SiteRepo.Get(workspace.SiteId)
	interps, _ := s.InterpreterService.List()
	mp, _ := s.InterpreterService.GetMap(interps)

	conf := configUtils.ReadFromFile(workspace.Path)
	if conf.Language == "" {
		conf.Language = commConsts.LanguageZh
	}
	if forceUpdate || conf.Url == "" {
		conf.Url = site.Url
	}
	if forceUpdate || conf.Username == "" {
		conf.Username = site.Username
	}
	if forceUpdate || conf.Password == "" {
		conf.Password = site.Password
	}

	if commonUtils.IsWin() {
		if forceUpdate || conf.Javascript == "" {
			conf.Javascript = mp["javascript"]
		}
		if forceUpdate || conf.Lua == "" {
			conf.Lua = mp["lua"]
		}
		if forceUpdate || conf.Perl == "" {
			conf.Perl = mp["perl"]
		}
		if forceUpdate || conf.Php == "" {
			conf.Php = mp["php"]
		}
		if forceUpdate || conf.Python == "" {
			conf.Python = mp["python"]
		}
		if forceUpdate || conf.Ruby == "" {
			conf.Ruby = mp["ruby"]
		}
		if forceUpdate || conf.Tcl == "" {
			conf.Tcl = mp["tcl"]
		}
		if forceUpdate || conf.Autoit == "" {
			conf.Autoit = mp["autoit"]
		}
	}

	configUtils.SaveToFile(conf, workspace.Path)

	return
}
