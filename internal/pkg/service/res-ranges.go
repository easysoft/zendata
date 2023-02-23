package service

import (
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/jinzhu/copier"
)

type ResRangesService struct {
	ResService   *ResService   `inject:""`
	FieldService *FieldService `inject:""`
}

func (s *ResRangesService) GetResFromRanges(ranges domain.ResRanges) map[string][]interface{} {
	groupedValue := map[string][]interface{}{}

	for group, expression := range ranges.Ranges {
		fieldFromRanges := s.ConvertRangesToField(ranges, expression)

		s.FieldService.Generate(&fieldFromRanges, true)

		groupedValue[group] = fieldFromRanges.Values
	}

	return groupedValue
}

func (s *ResRangesService) ConvertRangesToField(ranges domain.ResRanges, expression string) (field domain.DefField) {
	field.Note = "ConvertedFromRanges"

	copier.Copy(&field, ranges)
	field.Range = expression

	return field
}
