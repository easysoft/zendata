package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type ReferRepo struct {
	db *gorm.DB
}

func (r *ReferRepo) CreateDefault(ownerId uint, ownerType string) (err error) {
	refer := &model.ZdRefer{OwnerID: ownerId, OwnerType: ownerType}
	err = r.db.Create(&refer).Error

	return
}

func (r *ReferRepo) Create(refer *model.ZdRefer) (err error) {
	err = r.db.Create(&refer).Error
	return
}

func (r *ReferRepo) Get(fieldId uint, ownerType string) (refer model.ZdRefer, err error) {
	err = r.db.Where("ownerID=? AND ownerType=?", fieldId, ownerType).First(&refer).Error
	return
}

func (r *ReferRepo) Save(ref *model.ZdRefer) (err error) {
	err = r.db.Save(ref).Error
	return
}

func NewReferRepo(db *gorm.DB) *ReferRepo {
	return &ReferRepo{db: db}
}
