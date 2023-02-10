package gen

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	valueGen "github.com/easysoft/zendata/internal/pkg/gen/value"
)

func CreateField(field *domain.DefField) domain.FieldWithValues {
	fieldWithValue := domain.FieldWithValues{}

	if field.Type == "" { // set default
		field.Type = consts.FieldTypeList
	}
	if field.Length > 0 {
		field.Length = field.Length - len(field.Prefix) - len(field.Postfix)
		if field.Length < 0 {
			field.Length = 0
		}
	}

	if field.Type == consts.FieldTypeList {
		CreateListField(field, &fieldWithValue)
	} else if field.Type == consts.FieldTypeTimestamp {
		valueGen.CreateTimestampField(field, &fieldWithValue)
	} else if field.Type == consts.FieldTypeUlid {
		valueGen.CreateUlidField(field, &fieldWithValue)
	} else if field.Type == consts.FieldTypeArticle {
		CreateArticleField(field, &fieldWithValue)
	}

	return fieldWithValue
}
