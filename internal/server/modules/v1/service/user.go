package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type UserService struct {
	UserRepo *repo.UserRepo `inject:""`
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Paginate(req serverDomain.UserReqPaginate) (domain.PageData, error) {
	return s.UserRepo.Paginate(req)
}

// getRoles
func (s *UserService) getRoles(users ...*serverDomain.UserResponse) {
	s.UserRepo.GetRoles(users...)
}

func (s *UserService) FindByUserName(username string, ids ...uint) (serverDomain.UserResponse, error) {
	return s.UserRepo.FindByUserName(username, ids...)
}

func (s *UserService) FindPasswordByUserName(username string, ids ...uint) (serverDomain.LoginResponse, error) {
	return s.UserRepo.FindPasswordByUserName(username, ids...)
}

func (s *UserService) Create(req serverDomain.UserRequest) (uint, error) {
	return s.UserRepo.Create(req)
}

func (s *UserService) Update(id uint, req serverDomain.UserRequest) error {
	return s.UserRepo.Update(id, req)
}

func (s *UserService) IsAdminUser(id uint) (bool, error) {
	return s.UserRepo.IsAdminUser(id)
}

func (s *UserService) FindById(id uint) (serverDomain.UserResponse, error) {
	return s.UserRepo.FindById(id)
}

func (s *UserService) DeleteById(id uint) error {
	return s.UserRepo.DeleteById(id)
}

// AddRoleForUser add roles for user
func (s *UserService) AddRoleForUser(user *model.SysUser) error {
	return s.UserRepo.AddRoleForUser(user)
}

// DelToken 删除token
func (s *UserService) DelToken(token string) error {
	return s.UserRepo.DelToken(token)
}

// CleanToken 清空 token
func (s *UserService) CleanToken(authorityType int, userId string) error {
	return s.UserRepo.CleanToken(authorityType, userId)
}

func (s *UserService) UpdateAvatar(id uint, avatar string) error {
	return s.UserRepo.UpdateAvatar(id, avatar)
}
