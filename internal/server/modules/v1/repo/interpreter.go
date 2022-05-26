package repo

import (
	"errors"
	"fmt"

	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InterpreterRepo struct {
	DB *gorm.DB `inject:""`
}

func NewInterpreterRepo() *InterpreterRepo {
	return &InterpreterRepo{}
}

func (r *InterpreterRepo) List() (pos []model.Interpreter, err error) {
	db := r.DB.Model(&model.Interpreter{}).Where("NOT deleted")
	err = db.Find(&pos).Error
	fmt.Println(err)

	return
}

func (r *InterpreterRepo) Get(id uint) (po model.Interpreter, err error) {
	err = r.DB.Model(&model.Interpreter{}).
		Where("id = ?", id).
		Where("NOT deleted").
		First(&po).Error
	if err != nil {
		logUtils.Errorf(color.RedString("find interpreter by id failed, error: %s.", err.Error()))
		return
	}

	return
}

func (r *InterpreterRepo) Create(interpreter model.Interpreter) (id uint, err error) {
	po, err := r.FindDuplicate(interpreter.Lang, 0)
	if po.ID != 0 {
		return 0, errors.New(fmt.Sprintf("%s interpreter already exist.", interpreter.Lang))
	}

	err = r.DB.Model(&model.Interpreter{}).Create(&interpreter).Error
	if err != nil {
		logUtils.Errorf(color.RedString("create interpreter failed, error: %s.", err.Error()))
		return 0, err
	}

	id = interpreter.ID

	return
}

func (r *InterpreterRepo) Update(interpreter model.Interpreter) error {
	po, err := r.FindDuplicate(interpreter.Lang, interpreter.ID)
	if po.ID != 0 {
		return errors.New(fmt.Sprintf("%s运行环境已存在", interpreter.Lang))
	}

	err = r.DB.Model(&model.Interpreter{}).Where("id = ?", interpreter.ID).Updates(&interpreter).Error
	if err != nil {
		logUtils.Errorf(color.RedString("update interpreter failed, error: %s.", err.Error()))
		return err
	}

	return nil
}

func (r *InterpreterRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.Interpreter{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete interpreter by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *InterpreterRepo) FindDuplicate(lang string, id uint) (po model.Interpreter, err error) {
	db := r.DB.Model(&model.Interpreter{}).
		Where("NOT deleted").
		Where("lang = ?", lang)

	if id != 0 {
		db.Where("id != ?", id)
	}
	err = db.First(&po).Error

	return
}
