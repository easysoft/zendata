package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type ReferRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *ReferRepo) CreateDefault(ownerId uint, ownerType string) (err error) {
	refer := &model.ZdRefer{OwnerID: ownerId, OwnerType: ownerType}
	err = r.DB.Create(&refer).Error

	return
}

func (r *ReferRepo) Create(refer *model.ZdRefer) (err error) {
	err = r.DB.Create(&refer).Error
	return
}

func (r *ReferRepo) GetByOwnerId(fieldId uint) (refer model.ZdRefer, err error) {
	err = r.DB.Where("ownerID=?", fieldId).First(&refer).Error
	return
}
func (r *ReferRepo) GetByOwnerIdAndType(fieldId uint, ownerType string) (refer model.ZdRefer, err error) {
	err = r.DB.Where("ownerID=? AND ownerType=?", fieldId, ownerType).First(&refer).Error
	return
}

func (r *ReferRepo) Save(ref *model.ZdRefer) (err error) {
	err = r.DB.Save(ref).Error
	return
}

func NewReferRepo(db *gorm.DB) *ReferRepo {
	return &ReferRepo{DB: db}
}
