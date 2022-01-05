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

type ProjectRepo struct {
	DB *gorm.DB `inject:""`
}

func NewProjectRepo() *ProjectRepo {
	return &ProjectRepo{}
}

func (r *ProjectRepo) Paginate(req serverDomain.ProjectReqPaginate) (data domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.Project{}).Where("NOT deleted")

	if req.Keywords != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}
	if req.Enabled != "" {
		db = db.Where("disabled = ?", commonUtils.IsDisable(req.Enabled))
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("count project error", zap.String("error:", err.Error()))
		return
	}

	pos := make([]*model.Project, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&pos).Error
	if err != nil {
		logUtils.Errorf("query project error", zap.String("error:", err.Error()))
		return
	}

	data.Populate(pos, count, req.Page, req.PageSize)

	return
}

func (r *ProjectRepo) FindById(id uint) (po model.Project, err error) {
	err = r.DB.Model(&model.Project{}).Where("id = ?", id).First(&po).Error
	if err != nil {
		logUtils.Errorf("find project by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *ProjectRepo) Create(project model.Project) (id uint, err error) {
	if _, err := r.FindByName(project.Name); !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("%d", domain.BizErrNameExist.Code)
	}

	err = r.DB.Model(&model.Project{}).Create(&project).Error
	if err != nil {
		logUtils.Errorf("add project error", zap.String("error:", err.Error()))
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

func (r *ProjectRepo) BatchDelete(id uint) (err error) {
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

func (r *ProjectRepo) DeleteById(id uint, tx *gorm.DB) (err error) {
	err = tx.Model(&model.Project{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete project by id error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *ProjectRepo) DeleteChildren(ids []int, tx *gorm.DB) (err error) {
	err = tx.Model(&model.Project{}).Where("id IN (?)", ids).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("batch delete project error", zap.String("error:", err.Error()))
		return err
	}

	return nil
}

func (r *ProjectRepo) GetChildrenIds(id uint) (ids []int, err error) {
	tmpl := `
		WITH RECURSIVE project AS (
			SELECT * FROM biz_project WHERE id = %d
			UNION ALL
			SELECT child.* FROM biz_project child, project WHERE child.parent_id = project.id
		)
		SELECT id FROM project WHERE id != %d
    `
	sql := fmt.Sprintf(tmpl, id, id)
	err = r.DB.Raw(sql).Scan(&ids).Error
	if err != nil {
		logUtils.Errorf("get children project error", zap.String("error:", err.Error()))
		return
	}

	return
}

func (r *ProjectRepo) FindByName(name string, ids ...uint) (po model.Project, err error) {
	db := r.DB.Model(&model.Project{}).Where("name = ?", name)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err = db.First(&po).Error
	if err != nil {
		logUtils.Errorf("find project by name error", zap.String("name:", name), zap.Uints("ids:", ids), zap.String("error:", err.Error()))
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
		First(&currProject).Error

	return
}

func (r *ProjectRepo) RemoveDefaultTag() (err error) {
	err = r.DB.Model(&model.Project{}).
		Where("true").
		Update("is_default", false).Error

	return err
}
