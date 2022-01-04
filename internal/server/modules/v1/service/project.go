package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
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

func (s *ProjectService) Create(project model.Project) (uint, error) {
	return s.ProjectRepo.Create(project)
}

func (s *ProjectService) Update(id uint, project model.Project) error {
	return s.ProjectRepo.Update(id, project)
}

func (s *ProjectService) DeleteById(id uint) error {
	return s.ProjectRepo.BatchDelete(id)
}

func (s *ProjectService) GetByUser() (projects []model.Project, currProject model.Project, err error) {
	projects, err = s.ProjectRepo.ListProjectByUser()
	currProject, err = s.ProjectRepo.GetCurrProjectByUser()

	return
}
