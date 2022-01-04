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

type TestSetRepo struct {
	DB *gorm.DB `inject:""`
}

func NewTestSetRepo() *TestSetRepo {
	return &TestSetRepo{}
}

func (r *TestSetRepo) Paginate(req serverDomain.ReqPaginate) (data domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.TestSet{}).Where("NOT deleted")

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

	testSets := make([]*model.TestSet, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&testSets).Error
	if err != nil {
		logUtils.Errorf("query product error", zap.String("error:", err.Error()))
		return
	}

	data.Populate(testSets, count, req.Page, req.PageSize)

	return
}

func (r *TestSetRepo) FindById(id uint) (po model.TestSet, err error) {
	err = r.DB.Model(&model.TestSet{}).Where("id = ?", id).First(&po).Error
	if err != nil {
		logUtils.Errorf("find product by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *TestSetRepo) Create(po model.TestSet) (id uint, err error) {
	if _, err := r.FindByName(po.Name); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("%d", domain.BizErrNameExist.Code)
	}

	err = r.DB.Model(&model.TestSet{}).Create(&po).Error
	if err != nil {
		logUtils.Errorf("add product error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *TestSetRepo) Update(id uint, testSet model.TestSet) (err error) {
	err = r.DB.Model(&model.TestSet{}).Where("id = ?", id).Updates(&testSet).Error
	if err != nil {
		logUtils.Errorf("update product error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *TestSetRepo) DeleteById(id uint) (err error) {
	err = r.DB.Model(&model.TestSet{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete set by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *TestSetRepo) FindByName(name string, ids ...uint) (po model.TestSet, err error) {
	db := r.DB.Model(&model.TestSet{}).Where("name = ?", name)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err = db.First(&po).Error
	if err != nil {
		logUtils.Errorf("find product by name error %s", err.Error())
		return
	}

	return
}
