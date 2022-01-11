package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
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

func (s *ProjectService) GetByUser(currProjectPath string) (projects []model.Project, currProject model.Project, asset serverDomain.TestAsset, err error) {
	projects, err = s.ProjectRepo.ListProjectByUser()

	found := false
	for _, p := range projects {
		if p.Path == currProjectPath {
			found = true
			break
		}
	}

	if !found {
		if err != nil {
			logUtils.Errorf("db operation error %s", err.Error())
			return
		}

		name := fileUtils.GetDirName(currProjectPath)
		newLocalProject := model.Project{Path: currProjectPath, Name: name}

		_, err = s.ProjectRepo.Create(newLocalProject)
		if err != nil {
			logUtils.Errorf("db operation error %s", err.Error())
			return
		}

		projects, err = s.ProjectRepo.ListProjectByUser()
	}

	s.ProjectRepo.SetCurrProject(currProjectPath)

	currProject, err = s.ProjectRepo.GetCurrProjectByUser()
	if err != nil {
		logUtils.Errorf("db operation error %s", err.Error())
		return
	}

	asset, err = s.AssetService.LoadScripts(currProject.Path)

	return
}

func (s *ProjectService) SaveConfig(config commDomain.ProjectConfig) (err error) {
	currProject, err := s.ProjectRepo.GetCurrProjectByUser()
	if err != nil {
		return
	}

	serverConfig.SaveConfig(config, currProject.Path)

	return
}
