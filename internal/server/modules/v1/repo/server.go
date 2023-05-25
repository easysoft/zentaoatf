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

type ServerRepo struct {
	DB *gorm.DB `inject:""`
}

func NewServerRepo() *ServerRepo {
	return &ServerRepo{}
}

func (r *ServerRepo) List() (pos []model.Server, err error) {
	db := r.DB.Model(&model.Server{}).Where("NOT deleted")
	err = db.Find(&pos).Error

	return
}

func (r *ServerRepo) Get(id uint) (po model.Server, err error) {
	err = r.DB.Model(&model.Server{}).
		Where("id = ?", id).
		Where("NOT deleted").
		First(&po).Error
	if err != nil {
		logUtils.Errorf(color.RedString("find server by id failed, error: %s.", err.Error()))
		return
	}

	return
}

func (r *ServerRepo) Create(server model.Server) (id uint, err error) {
	po, err := r.FindDuplicate(server.Path, 0)
	if po.ID != 0 {
		return 0, errors.New(fmt.Sprintf("%s server already exist.", server.Path))
	}

	err = r.DB.Model(&model.Server{}).Create(&server).Error
	if err != nil {
		logUtils.Errorf(color.RedString("create server failed, error: %s.", err.Error()))
		return 0, err
	}

	id = server.ID

	return
}

func (r *ServerRepo) Update(server model.Server) error {
	po, err := r.FindDuplicate(server.Path, server.ID)
	if po.ID != 0 {
		return errors.New(fmt.Sprintf("%s远程服务器已存在", server.Path))
	}
	err = r.DB.Model(&model.Server{}).Where("id = ?", server.ID).Updates(&server).Error
	if err != nil {
		logUtils.Errorf(color.RedString("update server failed, error: %s.", err.Error()))
		return err
	}
	err = r.DB.Model(&model.Server{}).Where("id = ?", server.ID).
		Updates(map[string]interface{}{"is_default": server.IsDefault}).Error
	if err != nil {
		logUtils.Errorf("update server failed, error: %s.", err.Error())
		return err
	}

	return nil
}

func (r *ServerRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.Server{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete server by id error:%v", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *ServerRepo) FindDuplicate(path string, id uint) (po model.Server, err error) {
	db := r.DB.Model(&model.Server{}).
		Where("NOT deleted").
		Where("path = ?", path)

	if id != 0 {
		db.Where("id != ?", id)
	}
	err = db.First(&po).Error

	return
}
