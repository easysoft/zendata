package defServer

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/vari"
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
	}

	err = vari.GormDB.Save(field).Error
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
