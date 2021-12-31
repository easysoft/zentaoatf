package repo

import (
	"errors"
	"fmt"
	myZap "github.com/aaronchen2k/deeptest/internal/pkg/core/zap"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/core/casbin"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

type RoleRepo struct {
	DB *gorm.DB `inject:""`
}

func NewRoleRepo() *RoleRepo {
	return &RoleRepo{}
}

// Paginate
func (r *RoleRepo) Paginate(req serverDomain.RoleReqPaginate) (data domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.SysRole{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%s%%", req.Name))
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("获取角色总数错误", zap.String("错误:", err.Error()))
		return
	}

	var roles []*serverDomain.RoleResponse
	err = db.Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).Find(&roles).Error
	if err != nil {
		logUtils.Errorf("获取角色分页数据错误", zap.String("错误:", err.Error()))
		return
	}

	data.Result = roles
	data.Populate(roles, count, req.Page, req.PageSize)

	return
}

// FindByName
func (r *RoleRepo) FindByName(name string, ids ...uint) (serverDomain.RoleResponse, error) {
	role := serverDomain.RoleResponse{}
	db := r.DB.Model(&model.SysRole{}).Where("name = ?", name)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err := db.First(&role).Error
	if err != nil {
		logUtils.Errorf("根据名称查询角色错误", zap.String("名称:", name), zap.String("错误:", err.Error()))
		return role, err
	}
	return role, nil
}

// FindByName
func (r *RoleRepo) FindFirstAdminUser() (serverDomain.RoleResponse, error) {
	role := serverDomain.RoleResponse{}
	err :=r.DB.Model(&model.SysRole{}).Where("true").First(&role).Error

	if err != nil {
		logUtils.Errorf("管理员角色不存在，错误%s。", err.Error())
		return role, err
	}
	return role, nil
}

func (r *RoleRepo) Create(req serverDomain.RoleRequest) (uint, error) {
	role := model.SysRole{BaseRole: req.BaseRole}
	_, err := r.FindByName(req.Name)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logUtils.Errorf("角色名称已经被使用")
		return 0, err
	}

	err = r.DB.Create(&role).Error
	if err != nil {
		logUtils.Errorf("create data err ", zap.String("错误:", err.Error()))
		return 0, err
	}

	err = r.AddPermForRole(role.ID, req.Perms)
	if err != nil {
		logUtils.Errorf("添加权限到角色错误", zap.String("错误:", err.Error()))
		return 0, err
	}

	return role.ID, nil
}

func (r *RoleRepo) Update(id uint, req serverDomain.RoleRequest) error {
	if b, err := r.IsAdminRole(id); err != nil {
		return err
	} else if b {
		return errors.New("不能编辑超级管理员")
	}
	_, err := r.FindByName(req.Name, id)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logUtils.Errorf("角色名称已经被使用")
		return err
	}
	role := model.SysRole{BaseRole: req.BaseRole}
	err = r.DB.Model(&model.SysRole{}).Where("id = ?", id).Updates(&role).Error
	if err != nil {
		logUtils.Errorf("更新角色错误", zap.String("错误:", err.Error()))
		return err
	}
	err = r.AddPermForRole(role.ID, req.Perms)
	if err != nil {
		logUtils.Errorf("添加权限到角色错误", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

func (r *RoleRepo) IsAdminRole(id uint) (bool, error) {
	role, err := r.FindById(id)
	if err != nil {
		return false, err
	}
	return role.Name == serverConsts.AdminRoleName, nil
}

func (r *RoleRepo) FindById(id uint) (serverDomain.RoleResponse, error) {
	role := serverDomain.RoleResponse{}
	err := r.DB.Model(&model.SysRole{}).Where("id = ?", id).First(&role).Error
	if err != nil {
		logUtils.Errorf("根据id查询角色错误", zap.String("错误:", err.Error()))
		return role, err
	}
	return role, nil
}

func (r *RoleRepo) DeleteById(id uint) error {
	if b, err := r.IsAdminRole(id); err != nil {
		return err
	} else if b {
		return errors.New("不能删除超级管理员")
	}
	err := r.DB.Unscoped().Delete(&model.SysRole{}, id).Error
	if err != nil {
		logUtils.Errorf("删除角色错误", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

func (r *RoleRepo) FindInId(ids []string) ([]*serverDomain.RoleResponse, error) {
	roles := []*serverDomain.RoleResponse{}
	err := r.DB.Model(&model.SysRole{}).Where("id in ?", ids).Find(&roles).Error
	if err != nil {
		logUtils.Errorf("通过ids查询角色错误", zap.String("错误:", err.Error()))
		return nil, err
	}
	return roles, nil
}

// AddPermForRole
func (r *RoleRepo) AddPermForRole(id uint, perms [][]string) error {
	roleId := strconv.FormatUint(uint64(id), 10)
	oldPerms := casbin.GetPermissionsForUser(roleId)
	_, err := casbin.Instance().RemovePolicies(oldPerms)
	if err != nil {
		logUtils.Errorf("add policy err: %+v", zap.String("错误:", err.Error()))
		return err
	}

	if len(perms) == 0 {
		logUtils.Debug("没有权限")
		return nil
	}
	var newPerms [][]string
	for _, perm := range perms {
		newPerms = append(newPerms, append([]string{roleId}, perm...))
	}
	logUtils.Debugf("添加权限到角色", myZap.Strings("新权限", newPerms))
	_, err = casbin.Instance().AddPolicies(newPerms)
	if err != nil {
		logUtils.Errorf("add policy err: %+v", zap.String("错误:", err.Error()))
		return err
	}

	return nil
}

func (r *RoleRepo) GetRoleIds() ([]uint, error) {
	var roleIds []uint
	err := r.DB.Model(&model.SysRole{}).Select("id").Find(&roleIds).Error
	if err != nil {
		return roleIds, fmt.Errorf("获取角色ids错误 %w", err)
	}
	return roleIds, nil
}
