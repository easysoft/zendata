package service

import (
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/jinzhu/copier"
)

type ResRangesService struct {
	ResService   *ResService   `inject:""`
	FieldService *FieldService `inject:""`
}

func (s *ResRangesService) GetResFromRanges(ranges model.ResRanges) map[string][]interface{} {
	groupedValue := map[string][]interface{}{}

	for group, expression := range ranges.Ranges {
		fieldFromRanges := s.ConvertRangesToField(ranges, expression)

		s.FieldService.Generate(&fieldFromRanges, true)

		groupedValue[group] = fieldFromRanges.Values
	}

	return groupedValue
}

func (s *ResRangesService) ConvertRangesToField(ranges model.ResRanges, expression string) (field model.DefField) {
	field.Note = "ConvertedFromRanges"

	copier.Copy(&field, ranges)
	field.Range = expression

	return field
}
