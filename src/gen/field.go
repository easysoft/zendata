package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
)

func CreateField(field *model.DefField) model.FieldWithValues {
	fieldWithValue := model.FieldWithValues{}

	if field.Type == "" { // set default
		field.Type = constant.FieldTypeList
	}

	if field.Type == constant.FieldTypeList {
		CreateListField(field, &fieldWithValue)
	} else if field.Type == constant.FieldTypeTimestamp {
		CreateTimestampField(field, &fieldWithValue)
	} else if field.Type == constant.FieldTypeArticle {
		CreateArticleField(field, &fieldWithValue)
	}

	return fieldWithValue
}
