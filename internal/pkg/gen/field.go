package gen

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
)

func CreateField(field *model.DefField) model.FieldWithValues {
	fieldWithValue := model.FieldWithValues{}

	if field.Type == "" { // set default
		field.Type = constant.FieldTypeList
	}
	if field.Length > 0 {
		field.Length = field.Length - len(field.Prefix) - len(field.Postfix)
		if field.Length < 0 {
			field.Length = 0
		}
	}

	if field.Type == constant.FieldTypeList {
		CreateListField(field, &fieldWithValue)
	} else if field.Type == constant.FieldTypeTimestamp {
		CreateTimestampField(field, &fieldWithValue)
	} else if field.Type == constant.FieldTypeUlid {
		CreateUlidField(field, &fieldWithValue)
	} else if field.Type == constant.FieldTypeArticle {
		CreateArticleField(field, &fieldWithValue)
	}

	return fieldWithValue
}
