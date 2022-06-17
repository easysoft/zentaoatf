package repo

import (
	"fmt"
	"github.com/easysoft/zentaoatf/internal/server/core/dao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/pkg/domain"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	httpUtils "github.com/easysoft/zentaoatf/pkg/lib/http"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SiteRepo struct {
	DB *gorm.DB `inject:""`
}

func NewSiteRepo() *SiteRepo {
	return &SiteRepo{}
}

func (r *SiteRepo) Paginate(req serverDomain.ReqPaginate) (data domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.Site{}).Where("NOT deleted")

	if req.Keywords != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}
	if req.Enabled != "" {
		db = db.Where("disabled = ?", commonUtils.IsDisable(req.Enabled))
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("count site error", zap.String("error:", err.Error()))
		return
	}

	pos := make([]*model.Site, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&pos).Error
	if err != nil {
		logUtils.Errorf("query site error", zap.String("error:", err.Error()))
		return
	}

	data.Populate(pos, count, req.Page, req.PageSize)

	return
}

func (r *SiteRepo) Get(id uint) (po model.Site, err error) {
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

func (r *SiteRepo) Create(site *model.Site) (id uint, isDuplicate bool, err error) {
	site.Url = httpUtils.AddSepIfNeeded(site.Url)

	po, err := r.FindDuplicate(site.Name, site.Url, 0)
	if po.ID != 0 {
		isDuplicate = true
		return
	}

	err = r.DB.Model(&model.Site{}).Create(site).Error
	if err != nil {
		logUtils.Errorf(color.RedString("create site failed, error: %s.", err.Error()))
		return
	}

	id = site.ID

	return
}

func (r *SiteRepo) Update(site model.Site) (isDuplicate bool, err error) {
	po, err := r.FindDuplicate(site.Name, site.Url, site.ID)
	if po.ID != 0 {
		isDuplicate = true
		return
	}

	err = r.DB.Model(&model.Site{}).Where("id = ?", site.ID).Updates(&site).Error
	if err != nil {
		logUtils.Errorf(color.RedString("update site failed, error: %s.", err.Error()))
		return
	}

	return
}

func (r *SiteRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.Site{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete site by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *SiteRepo) FindDuplicate(name, url string, id uint) (po model.Site, err error) {
	db := r.DB.Model(&model.Site{}).
		Where("NOT deleted").
		Where("name = ?", name)

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
