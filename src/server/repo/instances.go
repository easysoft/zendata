package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type InstancesRepo struct {
	db *gorm.DB
}

func (r *InstancesRepo) List() (models []*model.ZdInstances, err error) {
	err = r.db.Select("id,title,name,folder,path").Where("true").Order("id ASC").Find(&models).Error
	return
}

func (r *InstancesRepo) Get(id uint) (model model.ZdInstances, err error) {
	err = r.db.Where("id=?", id).First(&model).Error
	return
}

func (r *InstancesRepo) Create(model *model.ZdInstances) (err error) {
	err = r.db.Create(model).Error
	return
}
func (r *InstancesRepo) Update(model *model.ZdInstances) (err error) {
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
	err = r.db.Where("instancesId=?", instancesId).Order("parentID ASC, ord ASC").Find(&items).Error
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

func (r *InstancesRepo) GetItemTree(instancesId int) (root model.ZdInstancesItem) {
	items, _ := r.GetItems(instancesId)

	root.ID = 0
	root.Field = "实例"
	r.makeTree(items, &root)

	return
}

func (r *InstancesRepo) makeTree(Data []*model.ZdInstancesItem, node *model.ZdInstancesItem) {
	children, _ := r.haveChild(Data, node)
	if children != nil {
		node.Fields = append(node.Fields, children[0:]...)
		for _, v := range children {
			_, has := r.haveChild(Data, v)
			if has {
				r.makeTree(Data, v)
			}
		}
	}
}
func (r *InstancesRepo) haveChild(Data []*model.ZdInstancesItem,
		node *model.ZdInstancesItem) (child []*model.ZdInstancesItem, yes bool) {
	for _, v := range Data {
		if v.ParentID == node.ID {
			child = append(child, v)
		}
	}
	if child != nil {
		yes = true
	}
	return
}

func NewInstancesRepo(db *gorm.DB) *InstancesRepo {
	return &InstancesRepo{db: db}
}
