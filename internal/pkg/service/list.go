package service

import (
	"strings"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
)

type ListService struct {
	TextService  *TextService  `inject:""`
	RangeService *RangeService `inject:""`
}

func (s *ListService) CreateListField(field *domain.DefField) {
	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			s.CreateListField(&child)
		}
	} else {
		s.CreateListFieldValues(field)
	}
}

func (s *ListService) CreateListFieldValues(field *domain.DefField) {
	if strings.Index(field.Range, ".txt") > -1 {
		s.TextService.CreateFieldValuesFromText(field)
	} else {
		s.RangeService.CreateFieldValuesFromRange(field)
	}
}

func (s *ListService) AppendValues(items *[]interface{}, val string, repeat int, total int) int {
	for round := 0; round < repeat; round++ {
		*items = append(*items, val)

		total++
		if total > consts.MaxNumb {
			break
		}
	}

	return total
}

func (s *ListService) AppendArrItems(items *[]interface{}, arr []string, total int, isRand bool) int {
	for i := 0; i < len(arr); i++ {
		idx := i
		if isRand {
			idx = commonUtils.RandNum(len(arr)) // should set random here too
		}

		*items = append(*items, arr[idx])

		total++
		if total > consts.MaxNumb {
			break
		}
	}

	return total
}
