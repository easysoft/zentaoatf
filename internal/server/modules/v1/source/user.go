package source

import (
	"github.com/gookit/color"
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"gorm.io/gorm"
)

type UserSource struct {
	UserRepo *repo.UserRepo `inject:""`
	RoleRepo *repo.RoleRepo `inject:""`
}

func NewUserSource() *UserSource {
	return &UserSource{}
}

func (s *UserSource) GetSources() ([]serverDomain.UserRequest, error) {
	roleIds, err := s.RoleRepo.GetRoleIds()
	if err != nil {
		return []serverDomain.UserRequest{}, err
	}
	users := []serverDomain.UserRequest{
		{
			BaseUser: model.BaseUser{
				Username: serverConsts.AdminUserName,
				Name:     "超级管理员",
				Intro:    "超级管理员",
				Avatar:   "images/avatar-m.svg",
			},
			Password: serverConsts.AdminUserPassword,
			RoleIds:  roleIds,
		},
	}
	return users, nil
}

func (s *UserSource) Init() error {
	return s.UserRepo.DB.Transaction(func(tx *gorm.DB) error {
		if tx.Model(&model.SysUser{}).Where("id IN ?", []int{1}).Find(&[]model.SysUser{}).RowsAffected == 1 {
			color.Danger.Printf("\n[Mysql] --> %s 表的初始数据已存在!", model.SysUser{}.TableName())
			return nil
		}
		sources, err := s.GetSources()
		if err != nil {
			return err
		}
		for _, source := range sources {
			if _, err := s.UserRepo.Create(source); err != nil { // 遇到错误时回滚事务
				return err
			}
		}
		color.Info.Printf("\n[Mysql] --> %s 表初始数据成功!", model.SysUser{}.TableName())
		return nil
	})
}
