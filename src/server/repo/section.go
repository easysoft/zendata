package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type SectionRepo struct {
	db *gorm.DB
}

func (r *SectionRepo) List(fieldId uint) (sections []*model.Section, err error) {
	err = r.db.Where("fieldID=?", fieldId).Find(&sections).Error
	return
}

func (r *SectionRepo) Get(id uint) (section *model.Section, err error) {
	err = r.db.Where("id=?", id).Find(&section).Error
	return
}

func (r *SectionRepo) Create(section *model.Section) (err error) {
	err = r.db.Create(&section).Error
	return
}

func (r *SectionRepo) Update(section *model.Section) (err error) {
	err = r.db.Save(&section).Error
	return
}

func (r *SectionRepo) Remove(id uint) (err error) {
	err = r.db.Where("id=?", id).Delete(&model.Section{}).Error
	return
}

func NewSectionRepo(db *gorm.DB) *SectionRepo {
	return &SectionRepo{db: db}
}
