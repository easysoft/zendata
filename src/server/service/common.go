package serverService

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"strings"
)

func genFieldFromZdField(treeNode model.ZdField, refer model.ZdRefer, field *model.DefField) {
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
	if field.Type == constant.FieldTypeList {
		field.Type = ""
	}
	if field.Mode == constant.ModeParallel {
		field.Mode = ""
	}

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

	// deal with refer
	if refer.Type != "" {
		logUtils.PrintTo(refer.Type)
		if refer.Type == "excel" {
			field.From = refer.File
			field.Select = refer.ColName
			field.Where = refer.Condition

		} else if refer.Type == constant.ResTypeRanges { // medium{2}
			field.From = refer.File
			field.Use = fmt.Sprintf("%s{%d}", refer.ColName, refer.Count)

		} else if refer.Type == constant.ResTypeInstances { // privateC{2}
			field.From = refer.File
			field.Use = fmt.Sprintf("%s{%d}", refer.ColName, refer.Count)

		} else if refer.Type == constant.ResTypeYaml { // dir/content.yaml{3}
			field.Range = fmt.Sprintf("%s{%d}", refer.File, refer.Count)

		} else if refer.Type == constant.ResTypeText { // dir/users.txt:2
			field.Range = fmt.Sprintf("%s:%d", refer.File, refer.Step)
		} else if refer.Type == constant.ResTypeConfig { // dir/users.txt:2
			field.Config = fmt.Sprintf("%s", refer.File)
		} else if refer.Type == constant.ResTypeValue { // dir/users.txt:2
			field.Value = fmt.Sprintf("%s", refer.Value)
		}
	}
}

func GetRelatedPathWithResDir(p string) (ret string) {
	rpl := vari.ZdPath + constant.ResDirYaml + constant.PthSep
	ret = strings.Replace(p, rpl, "", 1)

	rpl = vari.ZdPath + constant.ResDirUsers + constant.PthSep
	ret = strings.Replace(ret, rpl, "", 1)

	return
}

//func instancesItemToResInstForExport(item model.ZdInstancesItem) (inst model.ResInstancesItem) {
//	inst.Note = item.Note
//
//	for _, child := range item.Fields {
//		childField := model.DefField{}
//		instancesItemToResFieldForExport(*child, &childField)
//
//		inst.Fields = append(inst.Fields, childField)
//	}
//
//	if len(inst.Fields) == 0 {
//		inst.Fields = nil
//	}
//	if len(inst.Froms) == 0 {
//		inst.Froms = nil
//	}
//
//	return
//}
//func instancesItemToResFieldForExport(item model.ZdInstancesItem, field *model.DefField) {
//	for _, item := range item.Fields {
//		childField := model.DefField{}
//		instancesItemToResFieldForExport(*item, &childField)
//
//		field.Fields = append(field.Fields, childField)
//	}
//
//	for _, from := range item.Froms { // only one level
//		childField := model.DefField{}
//		genFieldFromZdInstancesItem(*from, &childField)
//
//		field.Froms = append(field.Froms, childField)
//	}
//
//	if len(field.Fields) == 0 {
//		field.Fields = nil
//	}
//	if len(field.Froms) == 0 {
//		field.Froms = nil
//	}
//}
//}
//func genFieldFromZdInstancesItem(item model.ZdInstancesItem, field *model.DefField) () {
//	field.Field = item.Field
//	field.Note = item.Note
//
//	field.Range = item.Range
//	field.Value = item.Exp
//	field.Prefix = item.Prefix
//	field.Postfix = item.Postfix
//	field.Loop = item.Loop
//	field.Loopfix = item.Loopfix
//	field.Format = item.Format
//	field.Type = item.Type
//	field.Mode = item.Mode
//	field.Length = item.Length
//	field.LeftPad = item.LeftPad
//	field.RightPad = item.RightPad
//	field.Rand = item.Rand
//
//	field.Config = item.Config
//	field.Use = item.Use
//	field.From = item.From
//	field.Select = item.Select
//	field.Where = item.Where
//	field.Limit = item.Limit
//}
