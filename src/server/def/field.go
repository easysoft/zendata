package defServer

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v3"
)

func GetDefFieldTree(defId uint) (root *model.Field, err error) {
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

	var def model.Def
	err = vari.GormDB.Where("id=?", field.DefID).First(&def).Error

	dataToYaml(&def)
	err = vari.GormDB.Model(&def).Where("id=?", def.ID).Update("yaml", def.Yaml).Error

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

func dataToYaml(def *model.Def) (str string) {
	root, err := GetDefFieldTree(def.ID)
	if err != nil {
		return
	}

	defData := model.DefData{}
	genDef(*def, &defData)

	for _, child := range root.Children { // ignore the root
		defField := model.DefField{}
		convertConfModel(*child, &defField)

		defData.Fields = append(defData.Fields, defField)
	}

	bytes, err := yaml.Marshal(defData)
	def.Yaml = string(bytes)

	return
}
func convertConfModel(treeNode model.Field, field *model.DefField) {
	genField(treeNode, field)

	for _, child := range treeNode.Children {
		defField := model.DefField{}
		convertConfModel(*child, &defField)

		field.Fields = append(field.Fields, defField)
	}

	for _, from := range treeNode.Froms { // only one level
		defField := model.DefField{}
		genField(*from, &defField)

		field.Froms = append(field.Froms, defField)
	}

	if len(field.Fields) == 0 {
		field.Fields = nil
	}
	if len(field.Froms) == 0 {
		field.Froms = nil
	}

	return
}
func genDef(def model.Def, data *model.DefData) () {
	data.Title = def.Title
	data.Desc = def.Desc
	data.Type = def.Type
}
func genField(treeNode model.Field, field *model.DefField) () {
	field.Field = treeNode.Field
	field.Note = treeNode.Note

	field.Range = treeNode.Range
	field.Value = treeNode.Exp
	field.Prefix = treeNode.Prefix
	field.Postfix = treeNode.Postfix
	field.Loop = treeNode.Loop
	field.Loopfix = treeNode.Loopfix
	field.Format = treeNode.Format
	field.Type = treeNode.Type
	field.Mode = treeNode.Mode
	field.Length = treeNode.Length
	field.LeftPad = treeNode.LeftPad
	field.RightPad = treeNode.RightPad
	field.Rand = treeNode.Rand

	field.Config = treeNode.Config
	field.Use = treeNode.Use
	field.From = treeNode.From
	field.Select = treeNode.Select
	field.Where = treeNode.Where
	field.Limit = treeNode.Limit
}
