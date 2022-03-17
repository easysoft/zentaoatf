package repo

import (
	"errors"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WorkspaceRepo struct {
	DB *gorm.DB `inject:""`
}

func NewWorkspaceRepo() *WorkspaceRepo {
	return &WorkspaceRepo{}
}

func (r *WorkspaceRepo) Paginate(req serverDomain.ReqPaginate) (data domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.Workspace{}).Where("NOT deleted")

	if req.Keywords != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}
	if req.Enabled != "" {
		db = db.Where("disabled = ?", commonUtils.IsDisable(req.Enabled))
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("count site error", zap.String("error:", err.Error()))
		return
	}

	pos := make([]*model.Workspace, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&pos).Error
	if err != nil {
		logUtils.Errorf("query site error", zap.String("error:", err.Error()))
		return
	}

	data.Populate(pos, count, req.Page, req.PageSize)

	return
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

func (r *WorkspaceRepo) Create(site model.Workspace) (id uint, err error) {
	po, err := r.FindDuplicate(site.Name, site.Path, 0)
	if po.ID != 0 {
		return 0, errors.New(fmt.Sprintf("工作目录%s（%s）已存在", site.Name, site.Path))
	}

	err = r.DB.Model(&model.Workspace{}).Create(&site).Error
	if err != nil {
		logUtils.Errorf(color.RedString("create site failed, error: %s.", err.Error()))
		return 0, err
	}

	id = site.ID

	return
}

func (r *WorkspaceRepo) Update(site model.Workspace) error {
	po, err := r.FindDuplicate(site.Name, site.Path, site.ID)
	if po.ID != 0 {
		return errors.New(fmt.Sprintf("工作目录%s(%s)已存在", site.Name, site.Path))
	}

	err = r.DB.Model(&model.Workspace{}).Where("id = ?", site.ID).Updates(&site).Error
	if err != nil {
		logUtils.Errorf(color.RedString("update site failed, error: %s.", err.Error()))
		return err
	}

	return nil
}

func (r *WorkspaceRepo) Delete(id uint) (err error) {
	err = r.DB.Where("id = ?", id).
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

func (r *WorkspaceRepo) FindDuplicate(name, url string, id uint) (po model.Workspace, err error) {
	db := r.DB.Model(&model.Workspace{}).
		Where("NOT deleted").
		Where("name = ? OR path = ?", name, url)

	if id != 0 {
		db.Where("id != ?", id)
	}
	err = db.First(&po).Error

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
