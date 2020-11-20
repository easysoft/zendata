package serverRepo

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/jinzhu/gorm"
)

type InstancesRepo struct {
	db *gorm.DB
}

func (r *InstancesRepo) ListAll() (models []*model.ZdInstances) {
	r.db.Select("id,title,name,folder,path,updatedAt").Find(&models)
	return
}

func (r *InstancesRepo) List(keywords string, page int) (models []*model.ZdInstances, total int, err error) {
	query := r.db.Select("id,title,name,folder,path").Order("id ASC")
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%")
	}
	if page > 0 {
		query = query.Offset((page-1) * constant.PageSize).Limit(constant.PageSize)
	}

	err = query.Find(&models).Error

	err = r.db.Model(&model.ZdInstances{}).Count(&total).Error

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
func (r *InstancesRepo) UpdateYaml(po model.ZdInstances) (err error) {
	err = r.db.Model(&model.ZdInstances{}).Where("id=?", po.ID).Update("yaml", po.Yaml).Error
	return
}

func (r *InstancesRepo) Remove(id uint) (err error) {
	inst := model.ZdInstances{}
	inst.ID = id

	err = r.db.Delete(inst).Error
	err = r.db.Where("instancesID = ?", id).Delete(&model.ZdInstancesItem{}).Error

	return
}

func (r *InstancesRepo) GetItems(instancesId uint) (items []*model.ZdInstancesItem, err error) {
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

func (r *InstancesRepo) GetItemTree(instancesId uint) (root model.ZdInstancesItem) {
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

func (r *InstancesRepo) GenInst(po model.ZdInstances, res *model.ResInsts) {
	res.Title = po.Title
	res.Desc = po.Desc
	res.Field = po.Field
}

func NewInstancesRepo(db *gorm.DB) *InstancesRepo {
	return &InstancesRepo{db: db}
}
