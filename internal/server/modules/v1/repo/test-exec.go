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
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TestExecRepo struct {
	DB          *gorm.DB     `inject:""`
	ProjectRepo *ProjectRepo `inject:""`
}

func NewTestExecRepo() *TestExecRepo {
	return &TestExecRepo{}
}

func (r *TestExecRepo) Paginate(req serverDomain.TestExecReqPaginate, projectPath string) (data domain.PageData, err error) {
	var count int64

	project, err := r.ProjectRepo.FindByPath(projectPath)
	if err != nil {
		return
	}

	db := r.DB.Model(&model.TestExec{}).
		Where("projectId = ?", project.ID).
		Where("NOT deleted")

	if req.Keywords != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}
	if req.Enabled != "" {
		db = db.Where("disabled = ?", commonUtils.IsDisable(req.Enabled))
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("count product error", zap.String("error:", err.Error()))
		return
	}

	testExecutions := make([]*model.TestExec, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&testExecutions).Error
	if err != nil {
		logUtils.Errorf("query product error", zap.String("error:", err.Error()))
		return
	}

	data.Populate(testExecutions, count, req.Page, req.PageSize)

	return
}

func (r *TestExecRepo) FindById(id uint) (model.TestExec, error) {
	product := model.TestExec{}
	err := r.DB.Model(&model.TestExec{}).Where("id = ?", id).First(&product).Error
	if err != nil {
		logUtils.Errorf("find product by id error", zap.String("error:", err.Error()))
		return product, err
	}

	return product, nil
}

func (r *TestExecRepo) Create(po model.TestExec) (id uint, err error) {
	if _, err := r.FindByName(po.Name); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("%d", domain.BizErrNameExist.Code)
	}

	err = r.DB.Model(&model.TestExec{}).Create(&po).Error
	if err != nil {
		logUtils.Errorf("add product error", zap.String("error:", err.Error()))
		return 0, err
	}

	id = po.ID

	return
}

func (r *TestExecRepo) Update(id uint, po model.TestExec) (err error) {
	err = r.DB.Model(&model.TestExec{}).Where("id = ?", id).Updates(&po).Error
	if err != nil {
		logUtils.Errorf("update product error", zap.String("error:", err.Error()))
		return err
	}

	return
}

func (r *TestExecRepo) DeleteById(id uint) (err error) {
	err = r.DB.Model(&model.TestExec{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete execution by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *TestExecRepo) FindByName(name string, ids ...uint) (po model.TestExec, err error) {
	db := r.DB.Model(&model.TestExec{}).Where("name = ?", name)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err = db.First(&po).Error
	if err != nil {
		logUtils.Errorf("find product by name error", zap.String("name:", name), zap.Uints("ids:", ids), zap.String("error:", err.Error()))
		return
	}

	return
}
