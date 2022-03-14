package repo

import (
	"errors"
	"fmt"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/fatih/color"
	"gorm.io/gorm"
)

type SiteRepo struct {
	DB *gorm.DB `inject:""`
}

func NewSiteRepo() *SiteRepo {
	return &SiteRepo{}
}

func (r *SiteRepo) List() (sites []model.Site, err error) {
	err = r.DB.Model(&model.Site{}).
		Where("NOT deleted").
		Find(&sites).Error

	return
}

func (r *SiteRepo) FindById(id uint) (po model.Site, err error) {
	err = r.DB.Model(&model.Site{}).
		Where("id = ?", id).
		Where("NOT deleted").
		First(&po).Error
	if err != nil {
		logUtils.Errorf(color.RedString("find site by id failed, error: %s.", err.Error()))
		return
	}

	return
}

func (r *SiteRepo) Create(project model.Site) (id uint, err error) {
	po, err := r.FindDuplicate(project.Name, project.Url, 0)
	if po.ID != 0 {
		return 0, errors.New(fmt.Sprintf("站点%s(%s)已存在", project.Name, project.Url))
	}

	err = r.DB.Model(&model.Site{}).Create(&project).Error
	if err != nil {
		logUtils.Errorf(color.RedString("create project failed, error: %s.", err.Error()))
		return 0, err
	}

	id = project.ID

	return
}

func (r *SiteRepo) Update(id uint, project model.Site) error {
	po, err := r.FindDuplicate(project.Name, project.Url, id)
	if po.ID != 0 {
		return errors.New(fmt.Sprintf("站点%s(%s)已存在", project.Name, project.Url))
	}

	err = r.DB.Model(&model.Site{}).Where("id = ?", id).Updates(&project).Error
	if err != nil {
		logUtils.Errorf(color.RedString("update project failed, error: %s.", err.Error()))
		return err
	}

	return nil
}

func (r *SiteRepo) DeleteByPath(pth string) (err error) {
	err = r.DB.Where("path = ?", pth).
		Delete(&model.Site{}).
		Error
	if err != nil {
		logUtils.Errorf(color.RedString("delete project failed, error: %s.", err.Error()))
		return
	}

	return
}

func (r *SiteRepo) FindDuplicate(name, url string, id uint) (po model.Site, err error) {
	db := r.DB.Model(&model.Site{}).
		Where("NOT deleted").
		Where("name = ? OR url = ?", name, url)

	if id != 0 {
		db.Where("id != ?", id)
	}
	err = db.First(&po).Error

	return
}

func (r *SiteRepo) GetCurrSiteByUser() (currSite model.Site, err error) {
	err = r.DB.Model(&model.Site{}).
		Where("is_default").
		Where("NOT deleted").
		First(&currSite).Error

	return
}

func (r *SiteRepo) RemoveDefaultTag() (err error) {
	err = r.DB.Model(&model.Site{}).
		Where("true").
		Update("is_default", false).Error

	return err
}

func (r *SiteRepo) SetCurrSite(id uint) (err error) {
	r.RemoveDefaultTag()

	if id == 0 {
		po := model.Site{}
		err := r.DB.Model(&model.Site{}).
			Where("NOT deleted").
			Order("id DESC").
			First(&po).Error
		if err == nil {
			id = po.ID
		}
	}

	err = r.DB.Model(&model.Site{}).
		Where("id = ?", id).
		Update("is_default", true).Error

	return err
}
