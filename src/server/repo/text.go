package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/jinzhu/gorm"
)

type TextRepo struct {
	db *gorm.DB
}

func (r *TextRepo) List(keywords string, page int) (models []*model.ZdText, total int, err error) {
	query := r.db.Select("id,title,folder,path").Order("id ASC")
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%")
	}
	if page > 0 {
		query = query.Offset((page-1) * constant.PageSize).Limit(constant.PageSize)
	}

	err = query.Find(&models).Error

	err = r.db.Model(&model.ZdText{}).Count(&total).Error

	return
}

func (r *TextRepo) Get(id uint) (model model.ZdText, err error) {
	err = r.db.Where("id=?", id).First(&model).Error
	return
}

func (r *TextRepo) Create(model *model.ZdText) (err error) {
	err = r.db.Create(model).Error
	return
}
func (r *TextRepo) Update(model *model.ZdText) (err error) {
	err = r.db.Save(model).Error
	return
}

func (r *TextRepo) Remove(id uint) (err error) {
	model := model.ZdText{}
	model.ID = id
	err = r.db.Delete(model).Error

	return
}

func NewTextRepo(db *gorm.DB) *TextRepo {
	return &TextRepo{db: db}
}
