package defServer

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/jinzhu/gorm"
)

func GetDefFieldTree(defId int) (root *model.Field, err error) {
	var fields []*model.Field

	err = vari.GormDB.Where("defID=?", defId).Order("parentID ASC ,ord ASC").Find(&fields).Error
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields")
	}

	root = fields[0]
	children := make([]*model.Field, 0)
	if len(fields) > 1 {
		children = fields[1:]
	}

	makeTree(children, root)
	return
}

func GetDefField(fieldId int) (field model.Field, err error) {
	err = vari.GormDB.Where("id=?", fieldId).First(&field).Error
	return
}

func SaveDefField(field *model.Field) (err error) {
	err = vari.GormDB.Save(field).Error
	return
}
func CreateDefField(defId, targetId uint, name string, mode string) (field *model.Field, err error) {
	field = &model.Field{Field: name, DefID: defId}
	if mode == "root" {
		field.DefID = defId
		field.ParentID = 0
	} else {
		var target model.Field

		err = vari.GormDB.Where("id=?", targetId).First(&target).Error
		field.DefID = target.DefID

		if mode == "child" {
			field.ParentID = target.ID
		} else {
			field.ParentID = target.ParentID
		}
		field.Ord = getMaxOrder(field.ParentID)
	}

	err = vari.GormDB.Save(field).Error
	return
}

func RemoveDefField(id int) (defId int, err error) {
	var field model.Field

	err = vari.GormDB.Where("id=?", id).First(&field).Error
	defId = int(field.DefID)
	err = deleteFieldAndChildren(field.DefID, field.ID)
	return
}

func MoveDefField(srcId, targetId uint, mode string) (defId int, srcField model.Field, err error) {
	var targetField model.Field
	err = vari.GormDB.Where("id=?", srcId).First(&srcField).Error
	if err == gorm.ErrRecordNotFound {
		return
	}
	err = vari.GormDB.Where("id=?", targetId).First(&targetField).Error
	if err == gorm.ErrRecordNotFound {
		return
	}
	defId = int(srcField.DefID)

	if "0" == mode {
		srcField.ParentID = targetId
		srcField.Ord = getMaxOrder(srcField.ParentID)
	} else if "-1" == mode {
		err = addOrderForTargetAndNextCases(srcField.ID, targetField.Ord, targetField.ParentID)
		if err != nil {
			return
		}

		srcField.ParentID = targetField.ParentID
		srcField.Ord = targetField.Ord
	} else if "1" == mode {
		err = addOrderForNextCases(srcField.ID, targetField.Ord, targetField.ParentID)
		if err != nil {
			return
		}

		srcField.ParentID = targetField.ParentID
		srcField.Ord = targetField.Ord + 1
	}

	sql := fmt.Sprintf(`update %s set ord = %d, parentID = %d where id=%d`,
		(&model.Field{}).TableName(), srcField.Ord, srcField.ParentID, srcField.ID)
	err = vari.GormDB.Exec(sql).Error

	return
}

func addOrderForTargetAndNextCases(srcID uint, targetOrder int, targetParentID uint) (err error) {
	sql := fmt.Sprintf(`update %s set ord = ord + 1 where ord >= %d and parentID = %d and id!=%d`,
		(&model.Field{}).TableName(), targetOrder, targetParentID, srcID)
	err = vari.GormDB.Exec(sql).Error

	return
}

func addOrderForNextCases(srcID uint, targetOrder int, targetParentID uint) (err error) {
	sql := fmt.Sprintf(`update %s set ord = ord + 1 where ord > %d and parentID = %d and id!=%d`,
		(&model.Field{}).TableName(), targetOrder, targetParentID, srcID)
	err = vari.GormDB.Exec(sql).Error

	return
}

func deleteFieldAndChildren(defId, fileId uint) (err error) {
	var children []*model.Field

	field := model.Field{}
	field.ID = fileId
	err = vari.GormDB.Delete(&field).Error

	if err == nil {
		err = vari.GormDB.Where("defID=? AND parentID=?", defId, fileId).Find(&children).Error
		for _, child := range children {
			deleteFieldAndChildren(child.DefID, child.ID)
		}
	}

	return
}

func makeTree(Data []*model.Field, node *model.Field) { //参数为父节点，添加父节点的子节点指针切片
	children, _ := haveChild(Data, node) //判断节点是否有子节点并返回
	if children != nil {
		node.Children = append(node.Children, children[0:]...) //添加子节点
		for _, v := range children {                           //查询子节点的子节点，并添加到子节点
			_, has := haveChild(Data, v)
			if has {
				makeTree(Data, v) //递归添加节点
			}
		}
	}
}

func haveChild(Data []*model.Field, node *model.Field) (child []*model.Field, yes bool) {
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

func getMaxOrder(parentId uint) (ord int) {
	var preChild model.Field
	err := vari.GormDB.
		Where("parentID=?", parentId).
		Order("ord DESC").Limit(1).
		First(&preChild).Error

	if err != nil {
		ord = 1
	}
	ord = preChild.Ord + 1

	return
}
