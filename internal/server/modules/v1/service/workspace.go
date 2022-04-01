package service

import (
	"errors"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type WorkspaceService struct {
	WorkspaceRepo *repo.WorkspaceRepo `inject:""`
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{}
}

func (s *WorkspaceService) Paginate(req serverDomain.WorkspaceReqPaginate) (ret domain.PageData, err error) {
	ret, err = s.WorkspaceRepo.Paginate(req)
	return
}

func (s *WorkspaceService) ListWorkspacesByProduct(siteId, productId int) (pos []model.Workspace, err error) {
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

	return
}
func (s *WorkspaceService) Update(workspace model.Workspace) (err error) {
	if !fileUtils.IsDir(workspace.Path) {
		err = errors.New(fmt.Sprintf("目录%s不存在", workspace.Path))
		return
	}

	err = s.WorkspaceRepo.Update(workspace)
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

func (s *WorkspaceService) ListWorkspaceByUser() (workspaces []model.Workspace, err error) {
	workspaces, err = s.WorkspaceRepo.ListWorkspaceByUser()

	return
}

func (s *WorkspaceService) GetByUser(currWorkspacePath string) (
	workspaces []model.Workspace, currWorkspace model.Workspace, currWorkspaceConfig commDomain.WorkspaceConf, scriptTree serverDomain.TestAsset, err error) {
	workspaces, err = s.WorkspaceRepo.ListWorkspaceByUser()

	found := false
	for _, p := range workspaces {
		if p.Path == currWorkspacePath {
			found = true
			break
		}
	}

	if !found {
		if err != nil {
			logUtils.Errorf("db operation error %s", err.Error())
			return
		}

		name := fileUtils.GetDirName(currWorkspacePath)
		newLocalWorkspace := model.Workspace{Path: currWorkspacePath, Name: name, Type: commConsts.ZTF}

		_, err = s.WorkspaceRepo.Create(newLocalWorkspace)
		if err != nil {
			logUtils.Errorf("db operation error %s", err.Error())
			return
		}

		workspaces, err = s.WorkspaceRepo.ListWorkspaceByUser()
	}

	s.WorkspaceRepo.SetCurrWorkspace(currWorkspacePath)

	currWorkspace, err = s.WorkspaceRepo.GetCurrWorkspaceByUser()
	if err != nil {
		logUtils.Errorf("db operation error %s", err.Error())
		return
	}

	if currWorkspace.Type == commConsts.ZTF {
		scriptTree, err = scriptUtils.LoadScriptTree(currWorkspace, nil)
	}

	currWorkspaceConfig = configUtils.ReadFromFile(currWorkspace.Path)
	currWorkspaceConfig.IsWin = commonUtils.IsWin()

	return
}

func (s *WorkspaceService) ListByProduct(siteId, productId int) (pos []model.Workspace, err error) {
	return s.WorkspaceRepo.ListByProduct(siteId, productId)
}
