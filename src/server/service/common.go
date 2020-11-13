package serverService

import "github.com/easysoft/zendata/src/model"

func convertToConfModel(treeNode model.Field, field *model.DefField) {
	genField(treeNode, field)

	for _, child := range treeNode.Children {
		defField := model.DefField{}
		convertToConfModel(*child, &defField)

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
