package repo

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"gorm.io/gorm"
	"time"
)

type JobRepo struct {
	DB *gorm.DB `inject:""`
}

func NewJobRepo() *JobRepo {
	return &JobRepo{}
}

func (r *JobRepo) Query() (pos []model.Job, err error) {
	err = r.DB.Model(&model.Job{}).
		Where("NOT deleted").
		Find(&pos).Error

	if err != nil {
		logUtils.Errorf("sql error %s", err.Error())
	}

	return
}

func (r *JobRepo) Get(id uint) (po model.Job, err error) {
	r.DB.Model(&model.Job{}).Where("id = ?", id).First(&po)

	return
}

func (r *JobRepo) Save(po *model.Job) (err error) {
	err = r.DB.Model(&model.Job{}).Create(&po).Error
	return
}

func (r *JobRepo) Update(po *model.Job) (err error) {
	err = r.DB.Model(&model.Job{}).Where("id = ?", po.ID).
		Session(&gorm.Session{FullSaveAssociations: true}).Updates(&po).Error
	return
}

func (r *JobRepo) UpdateStatus(id uint, status commConsts.JobStatus, isStart, isEnd bool) (err error) {

	updates := map[string]interface{}{"status": status}

	if isStart {
		updates["start_date"] = time.Now()
	}
	if isEnd {
		updates["end_date"] = time.Now()
	}

	err = r.DB.Model(&model.Job{}).Where("id = ?", id).
		Updates(updates).Error

	return
}

func (r *JobRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.Job{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true, "deleted_date": time.Now()}).Error

	return
}

func (r *JobRepo) SetFailed(po model.Job) (err error) {
	r.DB.Model(&model.Job{}).Where("id=?", po.ID).Updates(
		map[string]interface{}{"status": commConsts.JobFailed, "timeout_date": time.Now()})
	return
}

func (r *JobRepo) SetCanceled(po model.Job) (err error) {
	r.DB.Model(&model.Job{}).Where("id=?", po.ID).Updates(
		map[string]interface{}{"status": commConsts.JobCanceled, "cancel_date": time.Now()})
	return
}

func (r *JobRepo) AddRetry(po model.Job) (err error) {
	r.DB.Model(&model.Job{}).Where("id=?", po.ID).Updates(
		map[string]interface{}{"retry": gorm.Expr("retry + ?", 1)})
	return
}
