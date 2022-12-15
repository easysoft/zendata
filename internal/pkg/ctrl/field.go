package ctrl

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/internal/pkg/service"
)

type FieldCtrl struct {
	Field          *model.DefField
	FieldService   *service.FieldService
	ValueService   *service.ValueService
	ListService    *service.ListService    `inject:""`
	ArticleService *service.ArticleService `inject:""`
}

func (c *FieldCtrl) CreateField() {
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
		c.ListService.CreateListField(c.Field)
	} else if c.Field.Type == constant.FieldTypeArticle {
		c.ArticleService.CreateArticleField(c.Field)
	} else if c.Field.Type == constant.FieldTypeTimestamp {
		c.ValueService.CreateTimestampField(c.Field)
	} else if c.Field.Type == constant.FieldTypeUlid {
		c.ValueService.CreateUlidField(c.Field)
	}

	return
}
