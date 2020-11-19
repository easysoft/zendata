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
	inst := model.ZdInstances{}
	inst.ID = id

	err = r.db.Delete(inst).Error
	err = r.db.Where("instancesID = ?", id).Delete(&model.ZdInstancesItem{}).Error

	return
}

func (r *InstancesRepo) GetItems(instancesId int) (items []*model.ZdInstancesItem, err error) {
	err = r.db.Where("instancesId=?", instancesId).Find(&items).Error
	return
}
func (r *InstancesRepo) GetItem(itemId uint) (item model.ZdInstancesItem, err error) {
	err = r.db.Where("id=?", itemId).First(&item).Error
	return
}
func (r *InstancesRepo) SaveItem(item *model.ZdInstancesItem) (err error) {
	err = r.db.Save(item).Error
	return
}
func (r *InstancesRepo) RemoveItem(id uint) (err error) {
	item := model.ZdInstancesItem{}
	item.ID = id
	err = r.db.Delete(item).Error
	return
}
func (r *InstancesRepo) GetMaxOrder(instancesId int) (ord int) {
	var preChild model.ZdField
	err := r.db.
		Where("instancesID=?", instancesId).
		Order("ord DESC").Limit(1).
		First(&preChild).Error

	if err != nil {
		ord = 1
	}
	ord = preChild.Ord + 1

	return
}

func NewInstancesRepo(db *gorm.DB) *InstancesRepo {
	return &InstancesRepo{db: db}
}
