package ctrl

import (
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
