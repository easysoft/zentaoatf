package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type PermService struct {
	PermRepo *repo.PermRepo `inject:""`
}

func NewPermService() *PermService {
	return &PermService{}
}

// Paginate
func (s *PermService) Paginate(req serverDomain.PermReqPaginate) (data domain.PageData, err error) {
	return s.PermRepo.Paginate(req)
}

// FindByNameAndAct
// db *gorm.DB
// name 名称
// act 方法
// ids 当 ids 的 len = 1 ，排除次 id 数据
func (s *PermService) FindByNameAndAct(name, act string, ids ...uint) (serverDomain.PermResponse, error) {
	return s.PermRepo.FindByNameAndAct(name, act, ids...)
}

// Create
func (s *PermService) Create(req serverDomain.PermRequest) (uint, error) {
	return s.PermRepo.Create(req)
}

// CreatenInBatches
func (s *PermService) CreatenInBatches(perms []model.SysPerm) error {
	return s.PermRepo.CreateInBatches(perms)
}

// Update
func (s *PermService) Update(id uint, req serverDomain.PermRequest) error {
	return s.PermRepo.Update(id, req)
}

// checkNameAndAct
func (r *PermService) checkNameAndAct(req serverDomain.PermRequest, ids ...uint) bool {
	return r.PermRepo.CheckNameAndAct(req, ids...)
}

// FindById
func (s *PermService) FindById(id uint) (serverDomain.PermResponse, error) {
	return s.PermRepo.FindById(id)
}

// DeleteById
func (s *PermService) DeleteById(id uint) error {
	return s.PermRepo.DeleteById(id)
}

// GetPermsForRole
func (s *PermService) GetPermsForRole() ([][]string, error) {
	return s.PermRepo.GetPermsForRole()
}
