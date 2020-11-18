package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type SectionRepo struct {
	db *gorm.DB
}

func (r *SectionRepo) List(ownerId uint, ownerType string) (sections []*model.ZdSection, err error) {
	err = r.db.Where("ownerID=? AND ownerType=?", ownerId, ownerType).Find(&sections).Error
	return
}

func (r *SectionRepo) Get(id uint, ownerType string) (section model.ZdSection, err error) {
	err = r.db.Where("id=? AND ownerType=?", id, ownerType).First(&section).Error
	return
}

func (r *SectionRepo) Create(section *model.ZdSection) (err error) {
	err = r.db.Create(&section).Error
	return
}

func (r *SectionRepo) Update(section *model.ZdSection) (err error) {
	err = r.db.Save(&section).Error
	return
}

func (r *SectionRepo) Remove(id uint, ownerType string) (err error) {
	err = r.db.Where("id=? AND ownerType=?", id, ownerType).Delete(&model.ZdSection{}).Error
	return
}

func NewSectionRepo(db *gorm.DB) *SectionRepo {
	return &SectionRepo{db: db}
}
