package service

import (
	"github.com/easysoft/zendata/internal/pkg/domain"
)

type ResConfigService struct {
	FieldService *FieldService `inject:""`
	ExcelService *ExcelService `inject:""`
	RangeService *RangeService `inject:""`
}

func (s *ResConfigService) GetResForConfig(configRes domain.DefField) (groupedValue map[string][]interface{}) {
	groupedValue = map[string][]interface{}{}

	// config field is a standard field
	s.RangeService.DealwithFixRange(&configRes)
	s.FieldService.GenerateValuesForNoReferField(&configRes)

	groupedValue["all"] = configRes.Values

	return groupedValue
}
