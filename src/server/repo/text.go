package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type TextRepo struct {
	db *gorm.DB
}

func (r *TextRepo) List() (models []*model.ZdText, err error) {
	err = r.db.Where("true").Order("id ASC").Find(&models).Error
	return
}

func (r *TextRepo) Get(id uint) (model model.ZdText, err error) {
	err = r.db.Where("id=?", id).First(&model).Error
	return
}

func (r *TextRepo) Save(model *model.ZdText) (err error) {
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
