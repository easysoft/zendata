package serverRepo

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	"github.com/jinzhu/gorm"
)

type FieldRepo struct {
	db *gorm.DB
}

func (r *FieldRepo) GetDefFieldTree(defId uint) (root *model.ZdField, err error) {
	fields, err := r.ListByDef(defId)

	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields")
	}

	root = fields[0]
	children := make([]*model.ZdField, 0)
	if len(fields) > 1 {
		children = fields[1:]
	}

	r.makeTree(children, root)
	return
}

func (r *FieldRepo) CreateTreeNode(defId, targetId uint, name string, mode string) (field *model.ZdField, err error) {
	field = &model.ZdField{Field: name, DefID: defId}
	if mode == "root" {
		field.DefID = defId
		field.ParentID = 0
	} else {
		var target model.ZdField

		err = r.db.Where("id=?", targetId).First(&target).Error
		field.DefID = target.DefID

		if mode == "child" {
			field.ParentID = target.ID
		} else {
			field.ParentID = target.ParentID
		}
		field.Ord = r.GetMaxOrder(field.ParentID)
	}

	err = r.db.Save(field).Error
	return
}

func (r *FieldRepo) GetMaxOrder(parentId uint) (ord int) {
	var preChild model.ZdField
	err := r.db.
		Where("parentID=?", parentId).
		Order("ord DESC").Limit(1).
		First(&preChild).Error

	if err != nil {
		ord = 1
	}
	ord = preChild.Ord + 1

	return
}

func (r *FieldRepo) ListByDef(defId uint) (fields []*model.ZdField, err error) {
	err = r.db.Where("defID=?", defId).Order("parentID ASC ,ord ASC").Find(&fields).Error
	return
}

func (r *FieldRepo) Get(fieldId uint) (field model.ZdField, err error) {
	err = r.db.Where("id=?", fieldId).First(&field).Error
	return
}

func (r *FieldRepo) Save(field *model.ZdField) (err error) {
	err = r.db.Save(field).Error
	return
}

func (r *FieldRepo) makeTree(Data []*model.ZdField, node *model.ZdField) { //参数为父节点，添加父节点的子节点指针切片
	children, _ := r.haveChild(Data, node) //判断节点是否有子节点并返回
	if children != nil {
		node.Children = append(node.Children, children[0:]...) //添加子节点
		for _, v := range children {                           //查询子节点的子节点，并添加到子节点
			_, has := r.haveChild(Data, v)
			if has {
				r.makeTree(Data, v) //递归添加节点
			}
		}
	}
}

func (r *FieldRepo) haveChild(Data []*model.ZdField, node *model.ZdField) (child []*model.ZdField, yes bool) {
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

func (r *FieldRepo) Remove(id uint) (err error) {
	field := model.ZdField{}
	field.ID = id
	err = r.db.Delete(field).Error

	return
}

func (r *FieldRepo) GetChildren(defId, fieldId uint) (children []*model.ZdField, err error) {
	err = r.db.Where("defID=? AND parentID=?", defId, fieldId).Find(&children).Error
	return
}

func (r *FieldRepo) SetIsRange(fieldId uint, b bool) (err error) {
	err = r.db.Model(&model.ZdField{}).
		Where("id = ?", fieldId).Update("isRange", b).Error

	return
}

func (r *FieldRepo) AddOrderForTargetAndNextCases(srcID uint, targetOrder int, targetParentID uint) (err error) {
	sql := fmt.Sprintf(`update %s set ord = ord + 1 where ord >= %d and parentID = %d and id!=%d`,
		(&model.ZdField{}).TableName(), targetOrder, targetParentID, srcID)
	err = r.db.Exec(sql).Error

	return
}

func (r *FieldRepo) AddOrderForNextCases(srcID uint, targetOrder int, targetParentID uint) (err error) {
	sql := fmt.Sprintf(`update %s set ord = ord + 1 where ord > %d and parentID = %d and id!=%d`,
		(&model.ZdField{}).TableName(), targetOrder, targetParentID, srcID)
	err = r.db.Exec(sql).Error

	return
}

func (r *FieldRepo) UpdateOrdAndParent(field model.ZdField) (err error) {
	err = r.db.Model(&field).UpdateColumn(model.ZdField{Ord: field.Ord, ParentID: field.ParentID}).Error

	return
}

func NewFieldRepo(db *gorm.DB) *FieldRepo {
	return &FieldRepo{db: db}
}
