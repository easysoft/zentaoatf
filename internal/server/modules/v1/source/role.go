package source

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"github.com/gookit/color"
	"gorm.io/gorm"
)

type RoleSource struct {
	RoleRepo *repo.RoleRepo `inject:""`
	PermRepo *repo.PermRepo `inject:""`
}

func NewRoleSource() *RoleSource {
	return &RoleSource{}
}

func (s *RoleSource) GetSources() ([]serverDomain.RoleRequest, error) {
	perms, err := s.PermRepo.GetPermsForRole()
	if err != nil {
		return []serverDomain.RoleRequest{}, err
	}
	sources := []serverDomain.RoleRequest{
		{
			BaseRole: model.BaseRole{
				Name:        "admin",
				DisplayName: "超级管理员",
				Description: "超级管理员",
			},
			Perms: perms,
		},
	}
	return sources, err
}

func (s *RoleSource) Init() error {
	return s.RoleRepo.DB.Transaction(func(tx *gorm.DB) error {
		if tx.Model(&model.SysRole{}).Where("id IN ?", []int{1}).Find(&[]model.SysRole{}).RowsAffected == 1 {
			color.Danger.Printf("\n[Mysql] --> %s 表的初始数据已存在!", model.SysRole{}.TableName())
			return nil
		}
		sources, err := s.GetSources()
		if err != nil {
			return err
		}
		for _, source := range sources {
			if _, err := s.RoleRepo.Create(source); err != nil { // 遇到错误时回滚事务
				return err
			}
		}

		color.Info.Printf("\n[Mysql] --> %s 表初始数据成功!", model.SysRole{}.TableName())
		return nil
	})
}
