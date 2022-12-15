package service

import (
	"github.com/easysoft/zendata/internal/pkg/model"
	"strings"
)

type ListService struct {
	TextService  *TextService  `inject:""`
	RangeService *RangeService `inject:""`
}

func (s *ListService) CreateListField(field *model.DefField) {
	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			s.CreateListField(&child)
		}
	} else {
		s.CreateListFieldValues(field)
	}
}

func (s *ListService) CreateListFieldValues(field *model.DefField) {
	if strings.Index(field.Range, ".txt") > -1 {
		s.TextService.CreateFieldValuesFromText(field)
	} else {
		s.RangeService.CreateFieldValuesFromRange(field)
	}
}
