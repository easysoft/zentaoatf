package repo

import (
	"errors"
	"fmt"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProjectRepo struct {
	DB *gorm.DB `inject:""`
}

func NewProjectRepo() *ProjectRepo {
	return &ProjectRepo{}
}

func (r *ProjectRepo) FindById(id uint) (po model.Project, err error) {
	err = r.DB.Model(&model.Project{}).
		Where("id = ?", id).
		Where("NOT deleted").
		First(&po).Error
	if err != nil {
		logUtils.Errorf("find project by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *ProjectRepo) Create(project model.Project) (id uint, err error) {
	po, err := r.FindByName(project.Name)
	if po.ID != 0 {
		return 0, errors.New(fmt.Sprintf("名称为%s的项目已存在", project.Name))
	}

	err = r.DB.Model(&model.Project{}).Create(&project).Error
	if err != nil {
		logUtils.Errorf("create project error", zap.String("error:", err.Error()))
		return 0, err
	}

	id = project.ID

	return
}

func (r *ProjectRepo) Update(id uint, project model.Project) error {
	err := r.DB.Model(&model.Project{}).Where("id = ?", id).Updates(&project).Error
	if err != nil {
		logUtils.Errorf("update project error", zap.String("error:", err.Error()))
		return err
	}

	return nil
}

func (r *ProjectRepo) DeleteByPath(pth string) (err error) {
	err = r.DB.Where("path = ?", pth).
		Delete(&model.Project{}).
		Error
	if err != nil {
		logUtils.Errorf(color.RedString("delete project error", err.Error()))
		return
	}

	return
}

func (r *ProjectRepo) FindByName(name string, ids ...uint) (po model.Project, err error) {
	db := r.DB.Model(&model.Project{}).
		Where("NOT deleted").
		Where("name = ?", name)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err = db.First(&po).Error

	return
}

func (r *ProjectRepo) FindByPath(projectPath string) (po model.Project, err error) {
	db := r.DB.Model(&model.Project{}).Where("path = ?", projectPath)

	err = db.First(&po).Error
	if err != nil {
		logUtils.Errorf("find project by path error", err.Error())
		return
	}

	return
}

func (r *ProjectRepo) ListProjectByUser() (projects []model.Project, err error) {
	err = r.DB.Model(&model.Project{}).
		Where("NOT deleted").
		Find(&projects).Error

	return
}

func (r *ProjectRepo) GetCurrProjectByUser() (currProject model.Project, err error) {
	err = r.DB.Model(&model.Project{}).
		Where("is_default").
		Where("NOT deleted").
		First(&currProject).Error

	return
}

func (r *ProjectRepo) RemoveDefaultTag() (err error) {
	err = r.DB.Model(&model.Project{}).
		Where("true").
		Update("is_default", false).Error

	return err
}

func (r *ProjectRepo) SetCurrProject(pth string) (err error) {
	r.RemoveDefaultTag()

	if pth == "" {
		po := model.Project{}
		err := r.DB.Model(&model.Project{}).
			Where("NOT deleted").
			Order("id DESC").
			First(&po).Error
		if err == nil {
			pth = po.Path
		}
	}

	err = r.DB.Model(&model.Project{}).
		Where("path = ?", pth).
		Update("is_default", true).Error

	return err
}
