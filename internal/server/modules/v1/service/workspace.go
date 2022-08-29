package service

import (
	"errors"
	"fmt"

	"github.com/easysoft/zentaoatf/pkg/domain"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
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
	s.UpdateConfig(workspace, "all")

	return
}
func (s *WorkspaceService) Update(workspace model.Workspace) (err error) {
	if !fileUtils.IsDir(workspace.Path) {
		err = errors.New(fmt.Sprintf("目录%s不存在", workspace.Path))
		return
	}

	err = s.WorkspaceRepo.Update(workspace)
	s.UpdateConfig(workspace, "all")

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

func (s *WorkspaceService) DeleteByPath(path string, productId uint) (err error) {
	err = s.WorkspaceRepo.DeleteByPath(path, productId)
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
		if item.Type != commConsts.ZTF {
			continue
		}

		s.UpdateConfig(item, "interpreter")
	}
}

func (s *WorkspaceService) UpdateConfig(workspace model.Workspace, by string) (err error) {
	site, _ := s.SiteRepo.Get(workspace.SiteId)
	interps, _ := s.InterpreterService.List()
	mp, _ := s.InterpreterService.GetMap(interps)

	conf := configHelper.ReadFromFile(workspace.Path)
	if conf.Language == "" {
		conf.Language = commConsts.LanguageZh
	}

	if by == "all" || by == "site" {
		conf.Url = site.Url
		conf.Username = site.Username
		conf.Password = site.Password
	}

	if by == "all" || by == "interpreter" {
		conf.Javascript = mp["javascript"]
		conf.Lua = mp["lua"]
		conf.Perl = mp["perl"]
		conf.Php = mp["php"]
		conf.Python = mp["python"]
		conf.Go = mp["go"]
		conf.Ruby = mp["ruby"]
		conf.Tcl = mp["tcl"]
		conf.Autoit = mp["autoit"]
	}

	configHelper.SaveToFile(conf, workspace.Path)

	return
}
