package repo

import (
	"errors"
	"fmt"

	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"

	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProxyRepo struct {
	DB *gorm.DB `inject:""`
}

func NewProxyRepo() *ProxyRepo {
	return &ProxyRepo{}
}

func (r *ProxyRepo) List() (pos []model.Proxy, err error) {
	db := r.DB.Model(&model.Proxy{}).Where("NOT deleted")
	err = db.Find(&pos).Error

	return
}

func (r *ProxyRepo) Get(id uint) (po model.Proxy, err error) {
	err = r.DB.Model(&model.Proxy{}).
		Where("id = ?", id).
		Where("NOT deleted").
		First(&po).Error
	if err != nil {
		logUtils.Errorf(color.RedString("find proxy by id failed, error: %s.", err.Error()))
		return
	}

	return
}

func (r *ProxyRepo) Create(proxy model.Proxy) (id uint, err error) {
	po, err := r.FindDuplicate(proxy.Path, 0)
	if po.ID != 0 {
		return 0, errors.New(fmt.Sprintf("%s proxy already exist.", proxy.Path))
	}

	err = r.DB.Model(&model.Proxy{}).Create(&proxy).Error
	if err != nil {
		logUtils.Errorf(color.RedString("create proxy failed, error: %s.", err.Error()))
		return 0, err
	}

	id = proxy.ID

	return
}

func (r *ProxyRepo) Update(proxy model.Proxy) error {
	po, err := r.FindDuplicate(proxy.Path, proxy.ID)
	if po.ID != 0 {
		return errors.New(fmt.Sprintf("%s执行节点已存在", proxy.Path))
	}

	err = r.DB.Model(&model.Proxy{}).Where("id = ?", proxy.ID).Updates(&proxy).Error
	if err != nil {
		logUtils.Errorf(color.RedString("update proxy failed, error: %s.", err.Error()))
		return err
	}

	return nil
}

func (r *ProxyRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.Proxy{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete proxy by id error:%v", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *ProxyRepo) FindDuplicate(path string, id uint) (po model.Proxy, err error) {
	db := r.DB.Model(&model.Proxy{}).
		Where("NOT deleted").
		Where("path = ?", path)

	if id != 0 {
		db.Where("id != ?", id)
	}
	err = db.First(&po).Error

	return
}
