package serverRepo

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"gorm.io/gorm"
)

type InstancesRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *InstancesRepo) ListAll() (models []*model.ZdInstances) {
	r.DB.Select("id,title,referName,fileName,folder,path,updatedAt").Find(&models)
	return
}

func (r *InstancesRepo) List(keywords string, page int) (models []*model.ZdInstances, total int, err error) {
	query := r.DB.Select("id,title,referName,fileName,folder,path").Order("id ASC")
	if keywords != "" {
		query = query.Where("title LIKE ?", "%"+keywords+"%")
	}
	if page > 0 {
		query = query.Offset((page - 1) * consts.PageSize).Limit(consts.PageSize)
	}

	err = query.Find(&models).Error

	var total64 int64
	err = r.DB.Model(&model.ZdInstances{}).Count(&total64).Error

	total = int(total64)
	return
}

func (r *InstancesRepo) Get(id uint) (model model.ZdInstances, err error) {
	err = r.DB.Where("id=?", id).First(&model).Error
	return
}

func (r *InstancesRepo) Create(model *model.ZdInstances) (err error) {
	err = r.DB.Create(model).Error
	return
}
func (r *InstancesRepo) Update(model *model.ZdInstances) (err error) {
	err = r.DB.Save(model).Error
	return
}
func (r *InstancesRepo) UpdateYaml(po model.ZdInstances) (err error) {
	err = r.DB.Model(&model.ZdInstances{}).Where("id=?", po.ID).Update("yaml", po.Yaml).Error
	return
}

func (r *InstancesRepo) Remove(id uint) (err error) {
	inst := model.ZdInstances{}
	inst.ID = id

	err = r.DB.Delete(inst).Error
	err = r.DB.Where("instancesID = ?", id).Delete(&model.ZdInstancesItem{}).Error

	return
}

func (r *InstancesRepo) GetItems(instancesId uint) (items []*model.ZdInstancesItem, err error) {
	err = r.DB.Where("instancesId=?", instancesId).Order("parentID ASC, ord ASC").Find(&items).Error
	return
}
func (r *InstancesRepo) GetItem(itemId uint) (item model.ZdInstancesItem, err error) {
	err = r.DB.Where("id=?", itemId).First(&item).Error
	return
}
func (r *InstancesRepo) SaveItem(item *model.ZdInstancesItem) (err error) {
	err = r.DB.Save(item).Error
	return
}
func (r *InstancesRepo) UpdateItemRange(rang string, id uint) (err error) {
	err = r.DB.Model(&model.ZdInstancesItem{}).Where("id=?", id).Update("range", rang).Error

	return
}
func (r *InstancesRepo) SetIsRange(id uint, b bool) (err error) {
	err = r.DB.Model(&model.ZdInstancesItem{}).
		Where("id = ?", id).Update("isRange", b).Error

	return
}
func (r *InstancesRepo) RemoveItem(id uint) (err error) {
	item := model.ZdInstancesItem{}
	item.ID = id
	err = r.DB.Delete(item).Error
	return
}
func (r *InstancesRepo) GetMaxOrder(instancesId int) (ord int) {
	var preChild model.ZdField
	err := r.DB.
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

func (r *InstancesRepo) GenInst(po model.ZdInstances, res *model.ResInstances) {
	res.Title = po.Title
	res.Desc = po.Desc
}

func NewInstancesRepo(db *gorm.DB) *InstancesRepo {
	return &InstancesRepo{DB: db}
}
