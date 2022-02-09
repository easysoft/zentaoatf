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

type TestCaseRepo struct {
	DB *gorm.DB `inject:""`
}

func NewTestCaseRepo() *TestCaseRepo {
	return &TestCaseRepo{}
}

func (r *TestCaseRepo) Paginate(req serverDomain.TestCaseReqPaginate) (data domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.TestCase{}).Where("NOT deleted")

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

	testCases := make([]*model.TestCase, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&testCases).Error
	if err != nil {
		logUtils.Errorf("query product error", zap.String("error:", err.Error()))
		return
	}

	data.Populate(testCases, count, req.Page, req.PageSize)

	return
}

func (r *TestCaseRepo) FindById(id uint) (model.TestCase, error) {
	product := model.TestCase{}
	err := r.DB.Model(&model.TestCase{}).Where("id = ?", id).First(&product).Error
	if err != nil {
		logUtils.Errorf("find product by id error", zap.String("error:", err.Error()))
		return product, err
	}

	return product, nil
}

func (r *TestCaseRepo) Create(po model.TestCase) (id uint, err error) {
	if _, err := r.FindByName(po.Name); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("%d", domain.BizErrNameNotExist.Code)
	}

	err = r.DB.Model(&model.TestCase{}).Create(&po).Error
	if err != nil {
		logUtils.Errorf("add product error", zap.String("error:", err.Error()))
		return 0, err
	}

	id = po.ID

	return
}

func (r *TestCaseRepo) Update(id uint, po model.TestCase) (err error) {
	err = r.DB.Model(&model.TestCase{}).Where("id = ?", id).Updates(&po).Error
	if err != nil {
		logUtils.Errorf("update product error", zap.String("error:", err.Error()))
		return err
	}

	return
}

func (r *TestCaseRepo) DeleteById(id uint) (err error) {
	err = r.DB.Model(&model.TestCase{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete case by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *TestCaseRepo) FindByName(name string, ids ...uint) (po model.TestCase, err error) {
	db := r.DB.Model(&model.TestCase{}).Where("name = ?", name)
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
