package ctrl

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/internal/pkg/service"
)

type FieldCtrl struct {
	Field        *model.DefField
	FieldService *service.FieldService
	ValueService *service.ValueService
}

func (c *FieldCtrl) CreateField() {
	fieldWithValue := model.FieldWithValues{}

	if c.Field.Type == "" { // set default
		c.Field.Type = constant.FieldTypeList
	}
	if c.Field.Length > 0 {
		c.Field.Length = c.Field.Length - len(c.Field.Prefix) - len(c.Field.Postfix)
		if c.Field.Length < 0 {
			c.Field.Length = 0
		}
	}

	if c.Field.Type == constant.FieldTypeList {
		gen.CreateListField(c.Field, &fieldWithValue)
	} else if c.Field.Type == constant.FieldTypeArticle {
		gen.CreateArticleField(c.Field, &fieldWithValue)
	} else if c.Field.Type == constant.FieldTypeTimestamp {
		c.ValueService.CreateTimestampField(c.Field)
	} else if c.Field.Type == constant.FieldTypeUlid {
		c.ValueService.CreateUlidField(c.Field)
	}

	return
}
