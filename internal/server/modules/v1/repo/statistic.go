package repo

import (
	"github.com/easysoft/zentaoatf/internal/server/core/dao"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/facebookgo/inject"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type StatisticRepo struct {
	DB *gorm.DB `inject:""`
}

func NewStatisticRepo() *StatisticRepo {
	instance := &StatisticRepo{}
	if instance.DB == nil {
		instance.injectModule(instance)
	}
	return instance
}

func (r *StatisticRepo) injectModule(instance *StatisticRepo) {
	var g inject.Graph
	g.Logger = logUtils.LoggerStandard.Sugar()

	// inject objects
	if err := g.Provide(
		&inject.Object{Value: dao.GetDB()},
		&inject.Object{Value: instance},
	); err != nil {
		logUtils.Fatalf("provide usecase objects to the Graph: %v", err)
	}
	err := g.Populate()
	if err != nil {
		logUtils.Fatalf("populate the incomplete Objects: %v", err)
	}
}

func (r *StatisticRepo) Get(id uint) (po model.Statistic, err error) {
	err = r.DB.Model(&model.Statistic{}).
		Where("id = ?", id).
		Where("NOT deleted").
		First(&po).Error
	if err != nil {
		logUtils.Errorf(color.RedString("find statistics by id failed, error: %s.", err.Error()))
		return
	}

	return
}

func (r *StatisticRepo) GetByPath(path string) (po model.Statistic, err error) {
	err = r.DB.Model(&model.Statistic{}).
		Where("path = ?", path).
		Where("NOT deleted").
		First(&po).Error
	if err == gorm.ErrRecordNotFound {
		return po, nil
	}
	if err != nil {
		logUtils.Errorf(color.RedString("find statistics by id failed, error: %s.", err.Error()))
		return
	}

	return
}

func (r *StatisticRepo) Create(statistics *model.Statistic) (id uint, isDuplicate bool, err error) {

	po, err := r.FindDuplicate(statistics.Path, 0)
	if po.ID != 0 {
		isDuplicate = true
		return
	}

	err = r.DB.Model(&model.Statistic{}).Create(statistics).Error
	if err != nil {
		logUtils.Errorf(color.RedString("create statistics failed, error: %s.", err.Error()))
		return
	}

	id = statistics.ID

	return
}

func (r *StatisticRepo) Update(statistics model.Statistic) (isDuplicate bool, err error) {
	po, err := r.FindDuplicate(statistics.Path, statistics.ID)
	if po.ID != 0 {
		isDuplicate = true
		return
	}

	err = r.DB.Model(&model.Statistic{}).Where("id = ?", statistics.ID).Updates(&statistics).Error
	if err != nil {
		logUtils.Errorf(color.RedString("update statistics failed, error: %s.", err.Error()))
		return
	}

	return
}

func (r *StatisticRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.Statistic{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete statistics by id error:%v", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *StatisticRepo) FindDuplicate(path string, id uint) (po model.Statistic, err error) {
	db := r.DB.Model(&model.Statistic{}).
		Where("NOT deleted").
		Where("path = ?", path)

	if id != 0 {
		db.Where("id != ?", id)
	}
	err = db.First(&po).Error

	return
}

func (r *StatisticRepo) UpdateStatistic(path string, total, success, fail int, logPath string) (err error) {
	err = r.DB.Model(&model.Statistic{}).
		Where("path = ?", path).
		Updates(map[string]interface{}{
			"total":    total,
			"succ":     success,
			"fail":     fail,
			"FailLogs": logPath,
		}).Error

	return err
}
