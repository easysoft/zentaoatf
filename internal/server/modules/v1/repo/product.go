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

type ProductRepo struct {
	DB *gorm.DB `inject:""`
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{}
}

func (r *ProductRepo) Paginate(req serverDomain.ReqPaginate) (data domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.Product{}).Where("NOT deleted")

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

	pos := make([]*model.Product, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&pos).Error
	if err != nil {
		logUtils.Errorf("query product error %s", err.Error())
		return
	}

	data.Populate(pos, count, req.Page, req.PageSize)

	return
}

func (r *ProductRepo) FindById(id uint) (po model.Product, err error) {
	err = r.DB.Model(&model.Product{}).Where("id = ?", id).First(&po).Error
	if err != nil {
		logUtils.Errorf("find product by id error %s", err.Error())
		return
	}

	return
}

func (r *ProductRepo) Create(product model.Product) (id uint, err error) {
	if _, err := r.FindByName(product.Name); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("%d", domain.BizErrNameNotExist.Code)
	}

	err = r.DB.Model(&model.Product{}).Create(&product).Error
	if err != nil {
		logUtils.Errorf("add product error %s", err.Error())
		return 0, err
	}

	id = product.ID

	return
}

func (r *ProductRepo) Update(id uint, product model.Product) error {
	err := r.DB.Model(&model.Product{}).Where("id = ?", id).Updates(&product).Error
	if err != nil {
		logUtils.Errorf("update product error %s", err.Error())
		return err
	}

	return nil
}

func (r *ProductRepo) BatchDelete(id uint) (err error) {
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

func (r *ProductRepo) DeleteById(id uint, tx *gorm.DB) (err error) {
	err = tx.Model(&model.Product{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete product by id error %s", err.Error())
		return
	}

	return
}

func (r *ProductRepo) DeleteChildren(ids []int, tx *gorm.DB) (err error) {
	err = tx.Model(&model.Product{}).Where("id IN (?)", ids).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("batch delete product error %s", err.Error())
		return err
	}

	return nil
}

func (r *ProductRepo) GetChildrenIds(id uint) (ids []int, err error) {
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
		logUtils.Errorf("get children product error %s", err.Error())
		return
	}

	return
}

func (r *ProductRepo) FindByName(name string, ids ...uint) (po model.Product, err error) {
	db := r.DB.Model(&model.Product{}).Where("name = ?", name)
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
