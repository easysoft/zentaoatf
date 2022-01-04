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

type TestExecutionRepo struct {
	DB *gorm.DB `inject:""`
}

func NewTestExecutionRepo() *TestExecutionRepo {
	return &TestExecutionRepo{}
}

func (r *TestExecutionRepo) Paginate(req serverDomain.TestExecutionReqPaginate) (data domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.TestExecution{}).Where("NOT deleted")

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

	testExecutions := make([]*model.TestExecution, 0)

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

func (r *TestExecutionRepo) FindById(id uint) (model.TestExecution, error) {
	product := model.TestExecution{}
	err := r.DB.Model(&model.TestExecution{}).Where("id = ?", id).First(&product).Error
	if err != nil {
		logUtils.Errorf("find product by id error", zap.String("error:", err.Error()))
		return product, err
	}

	return product, nil
}

func (r *TestExecutionRepo) Create(po model.TestExecution) (id uint, err error) {
	if _, err := r.FindByName(po.Name); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("%d", domain.BizErrNameExist.Code)
	}

	err = r.DB.Model(&model.TestExecution{}).Create(&po).Error
	if err != nil {
		logUtils.Errorf("add product error", zap.String("error:", err.Error()))
		return 0, err
	}

	id = po.ID

	return
}

func (r *TestExecutionRepo) Update(id uint, po model.TestExecution) (err error) {
	err = r.DB.Model(&model.TestExecution{}).Where("id = ?", id).Updates(&po).Error
	if err != nil {
		logUtils.Errorf("update product error", zap.String("error:", err.Error()))
		return err
	}

	return
}

func (r *TestExecutionRepo) DeleteById(id uint) (err error) {
	err = r.DB.Model(&model.TestExecution{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete execution by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *TestExecutionRepo) FindByName(name string, ids ...uint) (po model.TestExecution, err error) {
	db := r.DB.Model(&model.TestExecution{}).Where("name = ?", name)
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
