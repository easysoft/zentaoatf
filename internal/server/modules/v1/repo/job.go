package repo

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/fatih/color"
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
	err = r.DB.
		Where("progress_status != ? && progress_status != ?",
			commConsts.ProgressCompleted, commConsts.ProgressTimeout).
		Where("NOT deleted").
		Find(&pos).Error

	if err != nil {
		logUtils.Errorf(color.RedString("list job failed: %s.", err.Error()))
		return
	}

	return
}

func (r *JobRepo) ListByProgressStatus(progress commConsts.ProgressStatus) (pos []model.Job, err error) {
	err = r.DB.
		Where("progress_status = ?", progress).
		Where("NOT deleted").
		Find(&pos).Error

	if err != nil {
		logUtils.Errorf(color.RedString("list job failed: %s.", err.Error()))
		return
	}

	return
}

func (r *JobRepo) Get(id uint) (po model.Job, err error) {
	err = r.DB.
		Where("id = ?", id).
		Where("NOT deleted").
		First(&po).Error

	if err != nil {
		logUtils.Errorf(color.RedString("get job by id failed: %s.", err.Error()))
		return
	}

	return
}

func (r *JobRepo) Create(job *model.Job) (err error) {
	err = r.DB.Model(&model.Job{}).Create(job).Error
	if err != nil {
		logUtils.Errorf(color.RedString("create job failed: %s.", err.Error()))
		return
	}

	return
}

func (r *JobRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.Job{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete job by id error: %s.", err.Error())
		return
	}

	return
}

func (r *JobRepo) SetTimeout(id uint) (err error) {
	err = r.DB.Model(&model.Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"progress_status": commConsts.ProgressTimeout,
			"timeout_time":    time.Now(),
		}).Error

	if err != nil {
		logUtils.Errorf("set job timeout error: %s.", err.Error())
		return
	}

	return
}

func (r *JobRepo) QueryForExec() (job model.Job, err error) {
	err = r.DB.
		Where("progress=?", commConsts.ProgressCreated).
		Order("priority ASC").
		First(&job).Error

	return
}

func (r *JobRepo) UpdateProgressStatus(id uint, progress commConsts.ProgressStatus) (err error) {
	err = r.DB.Model(&model.Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"progress_status": progress,
			"start_time":      time.Now(),
		}).Error

	if err != nil {
		logUtils.Errorf("set job progress error: %s.", err.Error())
		return
	}

	return
}
