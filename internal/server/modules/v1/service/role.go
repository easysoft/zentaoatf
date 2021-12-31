package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type RoleService struct {
	RoleRepo *repo.RoleRepo `inject:""`
}

func NewRoleService() *RoleService {
	return &RoleService{}
}

// Paginate
func (s *RoleService) Paginate(req serverDomain.RoleReqPaginate) (ret domain.PageData, err error) {
	return s.RoleRepo.Paginate(req)
}

// FindByName
func (s *RoleService) FindByName(name string, ids ...uint) (serverDomain.RoleResponse, error) {
	return s.RoleRepo.FindByName(name, ids...)
}

func (s *RoleService) Create(req serverDomain.RoleRequest) (uint, error) {
	return s.RoleRepo.Create(req)
}

func (s *RoleService) Update(id uint, req serverDomain.RoleRequest) error {
	return s.RoleRepo.Update(id, req)
}

func (s *RoleService) IsAdminRole(id uint) (bool, error) {
	return s.RoleRepo.IsAdminRole(id)
}

func (s *RoleService) FindById(id uint) (serverDomain.RoleResponse, error) {
	return s.RoleRepo.FindById(id)
}

func (s *RoleService) DeleteById(id uint) error {
	return s.RoleRepo.DeleteById(id)
}

func (s *RoleService) FindInId(ids []string) ([]*serverDomain.RoleResponse, error) {
	return s.RoleRepo.FindInId(ids)
}

// AddPermForRole
func (s *RoleService) AddPermForRole(id uint, perms [][]string) error {
	return s.RoleRepo.AddPermForRole(id, perms)
}

func (s *RoleService) GetRoleIds() ([]uint, error) {
	return s.RoleRepo.GetRoleIds()
}
