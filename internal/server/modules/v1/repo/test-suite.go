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

type TestSuiteRepo struct {
	DB *gorm.DB `inject:""`
}

func NewTestSuiteRepo() *TestSuiteRepo {
	return &TestSuiteRepo{}
}

func (r *TestSuiteRepo) Paginate(req serverDomain.ReqPaginate) (data domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.TestSuite{}).Where("NOT deleted")

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

	testSuites := make([]*model.TestSuite, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&testSuites).Error
	if err != nil {
		logUtils.Errorf("query product error", zap.String("error:", err.Error()))
		return
	}

	data.Populate(testSuites, count, req.Page, req.PageSize)

	return
}

func (r *TestSuiteRepo) FindById(id uint) (testSuite model.TestSuite, err error) {
	err = r.DB.Model(&model.TestSuite{}).Where("id = ?", id).First(&testSuite).Error
	if err != nil {
		logUtils.Errorf("find product by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *TestSuiteRepo) Create(testSuite model.TestSuite) (id uint, err error) {
	if _, err := r.FindByName(testSuite.Name); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("%d", domain.BizErrNameNotExist.Code)
	}

	err = r.DB.Model(&model.TestSuite{}).Create(&testSuite).Error
	if err != nil {
		logUtils.Errorf("add product error %s", err.Error())
		return 0, err
	}

	id = testSuite.ID

	return
}

func (r *TestSuiteRepo) Update(id uint, testSuite model.TestSuite) (err error) {
	err = r.DB.Model(&model.TestSuite{}).Where("id = ?", id).Updates(&testSuite).Error
	if err != nil {
		logUtils.Errorf("update product error %s", err.Error())
		return err
	}

	return nil
}

func (r *TestSuiteRepo) DeleteById(id uint) (err error) {
	err = r.DB.Model(&model.TestSuite{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete suite by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *TestSuiteRepo) FindByName(name string, ids ...uint) (testSuite model.TestSuite, err error) {
	db := r.DB.Model(&model.TestSuite{}).Where("name = ?", testSuite)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err = db.First(&testSuite).Error
	if err != nil {
		logUtils.Errorf("find product by name error %s", err.Error())
		return
	}

	return
}
