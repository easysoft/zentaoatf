package repo

import (
	"errors"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TestCaseRepo struct {
	DB       *gorm.DB  `inject:""`
	RoleRepo *RoleRepo `inject:""`
}

func NewTestCaseRepo() *TestCaseRepo {
	return &TestCaseRepo{}
}

func (r *TestCaseRepo) Paginate(req serverDomain.TestCaseReqPaginate) (data domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.TestCase{}).
		Where("NOT deleted")
	if req.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%s%%", req.Name))
	}
	if req.Category != "" {
		db = db.Where("category = ?", req.Category)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
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

func (r *TestCaseRepo) FindById(id uint) (serverDomain.TestCaseResponse, error) {
	product := serverDomain.TestCaseResponse{}
	err := r.DB.Model(&model.TestCase{}).Where("id = ?", id).First(&product).Error
	if err != nil {
		logUtils.Errorf("find product by id error", zap.String("error:", err.Error()))
		return product, err
	}

	return product, nil
}

func (r *TestCaseRepo) Create(req serverDomain.TestCaseRequest) (uint, error) {
	if _, err := r.FindByName(req.Name); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("%d", domain.BizErrNameExist.Code)
	}
	product := req.TestCase

	err := r.DB.Model(&model.TestCase{}).Create(&product).Error
	if err != nil {
		logUtils.Errorf("add product error", zap.String("error:", err.Error()))
		return 0, err
	}

	return product.ID, nil
}

func (r *TestCaseRepo) Update(id uint, req serverDomain.TestCaseRequest) error {
	product := req.TestCase
	err := r.DB.Model(&model.TestCase{}).Where("id = ?", id).Updates(&product).Error
	if err != nil {
		logUtils.Errorf("update product error", zap.String("error:", err.Error()))
		return err
	}

	return nil
}

func (r *TestCaseRepo) BatchDelete(id uint) (err error) {
	ids, err := r.GetChildrenIds(id)
	if err != nil {
		return err
	}

	r.DB.Transaction(func(tx *gorm.DB) (err error) {
		err = r.DeleteChildren(ids, tx)
		if err != nil {
			return
		}

		err = r.DeleteById(id, tx)
		if err != nil {
			return
		}

		return
	})

	return
}

func (r *TestCaseRepo) DeleteById(id uint, tx *gorm.DB) (err error) {
	err = tx.Model(&model.TestCase{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete product by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *TestCaseRepo) DeleteChildren(ids []int, tx *gorm.DB) (err error) {
	err = tx.Model(&model.TestCase{}).Where("id IN (?)", ids).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("batch delete product error", zap.String("error:", err.Error()))
		return err
	}

	return nil
}

func (r *TestCaseRepo) GetChildrenIds(id uint) (ids []int, err error) {
	tmpl := `
		WITH RECURSIVE product AS (
			SELECT * FROM biz_product WHERE id = %d
			UNION ALL
			SELECT child.* FROM biz_product child, product WHERE child.parent_id = product.id
		)
		SELECT id FROM product WHERE id != %d
    `
	sql := fmt.Sprintf(tmpl, id, id)
	err = r.DB.Raw(sql).Scan(&ids).Error
	if err != nil {
		logUtils.Errorf("get children product error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *TestCaseRepo) FindByName(productname string, ids ...uint) (serverDomain.TestCaseResponse, error) {
	product := serverDomain.TestCaseResponse{}
	db := r.DB.Model(&model.TestCase{}).Where("name = ?", productname)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err := db.First(&product).Error
	if err != nil {
		logUtils.Errorf("find product by name error", zap.String("name:", productname), zap.Uints("ids:", ids), zap.String("error:", err.Error()))
		return product, err
	}

	return product, nil
}