package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type InstancesRepo struct {
	db *gorm.DB
}

func (r *InstancesRepo) List() (models []*model.ZdInstances, err error) {
	err = r.db.Where("true").Order("id ASC").Find(&models).Error
	return
}

func (r *InstancesRepo) Get(id uint) (model model.ZdInstances, err error) {
	err = r.db.Where("id=?", id).First(&model).Error
	return
}

func (r *InstancesRepo) Save(model *model.ZdInstances) (err error) {
	err = r.db.Save(model).Error
	return
}

func (r *InstancesRepo) Remove(id uint) (err error) {
	model := model.ZdInstances{}
	model.ID = id
	err = r.db.Delete(model).Error

	return
}

func NewInstancesRepo(db *gorm.DB) *InstancesRepo {
	return &InstancesRepo{db: db}
}
