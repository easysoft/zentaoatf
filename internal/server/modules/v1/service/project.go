package service

import (
	"errors"
	"fmt"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/config"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/script"
)

type ProjectService struct {
	ProjectRepo *repo.ProjectRepo `inject:""`
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

func (s *ProjectService) Create(project model.Project) (id uint, err error) {
	po, _ := s.ProjectRepo.FindByPath(fileUtils.AddPathSepIfNeeded(project.Path))

	if po.ID != 0 {
		err = errors.New(fmt.Sprintf("路径为%s的项目已存在。", project.Path))
		return
	}

	id, _ = s.ProjectRepo.Create(project)
	return
}

func (s *ProjectService) Update(id uint, project model.Project) error {
	return s.ProjectRepo.Update(id, project)
}

func (s *ProjectService) DeleteById(id uint) error {
	return s.ProjectRepo.BatchDelete(id)
}

func (s *ProjectService) GetByUser(currProjectPath string) (
	projects []model.Project, currProject model.Project, currProjectConfig commDomain.ProjectConf, scriptTree serverDomain.TestAsset, err error) {
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

	scriptTree, err = scriptUtils.LoadScriptTree(currProject.Path)

	currProjectConfig = configUtils.ReadFromFile(currProject.Path)

	return
}
