package repo

import (
	"errors"
	"fmt"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/fatih/color"
	"gorm.io/gorm"
)

type WorkspaceRepo struct {
	DB *gorm.DB `inject:""`
}

func NewWorkspaceRepo() *WorkspaceRepo {
	return &WorkspaceRepo{}
}

func (r *WorkspaceRepo) FindById(id uint) (po model.Workspace, err error) {
	err = r.DB.Model(&model.Workspace{}).
		Where("id = ?", id).
		Where("NOT deleted").
		First(&po).Error
	if err != nil {
		logUtils.Errorf(color.RedString("find workspace by id failed, error: %s.", err.Error()))
		return
	}

	return
}

func (r *WorkspaceRepo) Create(workspace model.Workspace) (id uint, err error) {
	po, err := r.FindByName(workspace.Name)
	if po.ID != 0 {
		return 0, errors.New(fmt.Sprintf("名称为%s的项目已存在", workspace.Name))
	}

	err = r.DB.Model(&model.Workspace{}).Create(&workspace).Error
	if err != nil {
		logUtils.Errorf(color.RedString("create workspace failed, error: %s.", err.Error()))
		return 0, err
	}

	id = workspace.ID

	return
}

func (r *WorkspaceRepo) Update(id uint, workspace model.Workspace) error {
	err := r.DB.Model(&model.Workspace{}).Where("id = ?", id).Updates(&workspace).Error
	if err != nil {
		logUtils.Errorf(color.RedString("update workspace failed, error: %s.", err.Error()))
		return err
	}

	return nil
}

func (r *WorkspaceRepo) DeleteByPath(pth string) (err error) {
	err = r.DB.Where("path = ?", pth).
		Delete(&model.Workspace{}).
		Error
	if err != nil {
		logUtils.Errorf(color.RedString("delete workspace failed, error: %s.", err.Error()))
		return
	}

	return
}

func (r *WorkspaceRepo) FindByName(name string, ids ...uint) (po model.Workspace, err error) {
	db := r.DB.Model(&model.Workspace{}).
		Where("NOT deleted").
		Where("name = ?", name)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err = db.First(&po).Error

	return
}

func (r *WorkspaceRepo) FindByPath(workspacePath string) (po model.Workspace, err error) {
	db := r.DB.Model(&model.Workspace{}).Where("path = ?", workspacePath)

	err = db.First(&po).Error
	if err != nil {
		logUtils.Errorf("find workspace by path error", err.Error())
		return
	}

	return
}

func (r *WorkspaceRepo) ListWorkspaceByUser() (workspaces []model.Workspace, err error) {
	err = r.DB.Model(&model.Workspace{}).
		Where("NOT deleted").
		Find(&workspaces).Error

	return
}

func (r *WorkspaceRepo) GetCurrWorkspaceByUser() (currWorkspace model.Workspace, err error) {
	err = r.DB.Model(&model.Workspace{}).
		Where("is_default").
		Where("NOT deleted").
		First(&currWorkspace).Error

	return
}

func (r *WorkspaceRepo) RemoveDefaultTag() (err error) {
	err = r.DB.Model(&model.Workspace{}).
		Where("true").
		Update("is_default", false).Error

	return err
}

func (r *WorkspaceRepo) SetCurrWorkspace(pth string) (err error) {
	r.RemoveDefaultTag()

	if pth == "" {
		po := model.Workspace{}
		err := r.DB.Model(&model.Workspace{}).
			Where("NOT deleted").
			Order("id DESC").
			First(&po).Error
		if err == nil {
			pth = po.Path
		}
	}

	err = r.DB.Model(&model.Workspace{}).
		Where("path = ?", pth).
		Update("is_default", true).Error

	return err
}
