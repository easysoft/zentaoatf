package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type ProjectService struct {
	ProjectRepo  *repo.ProjectRepo `inject:""`
	AssetService *AssetService     `inject:""`
}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

func (s *ProjectService) Paginate(req serverDomain.ProjectReqPaginate) (ret domain.PageData, err error) {

	ret, err = s.ProjectRepo.Paginate(req)

	if err != nil {
		return
	}

	return
}

func (s *ProjectService) FindById(id uint) (model.Project, error) {
	return s.ProjectRepo.FindById(id)
}

func (s *ProjectService) Create(project model.Project) (uint, error) {
	return s.ProjectRepo.Create(project)
}

func (s *ProjectService) Update(id uint, project model.Project) error {
	return s.ProjectRepo.Update(id, project)
}

func (s *ProjectService) DeleteById(id uint) error {
	return s.ProjectRepo.BatchDelete(id)
}

func (s *ProjectService) GetByUser() (projects []model.Project, currProject model.Project, asset serverDomain.TestAsset, err error) {
	projects, err = s.ProjectRepo.ListProjectByUser()

	found := false
	for _, p := range projects {
		if p.Path == serverConfig.CONFIG.System.WorkDir {
			found = true
			break
		}
	}

	if !found {
		err = s.ProjectRepo.RemoveDefaultTag()
		if err != nil {
			logUtils.Errorf("db operation error %s", err.Error())
			return
		}

		pth := serverConfig.CONFIG.System.WorkDir
		name := fileUtils.GetDirName(pth)
		newLocalProject := model.Project{Path: pth, Name: name, IsDefault: true}
		_, err = s.ProjectRepo.Create(newLocalProject)
		if err != nil {
			logUtils.Errorf("db operation error %s", err.Error())
			return
		}

		projects, err = s.ProjectRepo.ListProjectByUser()
	}

	currProject, err = s.ProjectRepo.GetCurrProjectByUser()
	if err != nil {
		logUtils.Errorf("db operation error %s", err.Error())
		return
	}

	asset, err = s.AssetService.LoadScripts(currProject.Path)

	return
}
