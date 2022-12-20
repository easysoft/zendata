package service

import (
	"github.com/easysoft/zendata/internal/pkg/model"
)

type ResConfigService struct {
	FieldService *FieldService `inject:""`
	ExcelService *ExcelService `inject:""`
}

func (s *ResConfigService) GetResForConfig(configRes model.DefField) (groupedValue map[string][]interface{}) {
	groupedValue = map[string][]interface{}{}

	// config field is a standard field
	s.FieldService.GenerateValuesForNoReferField(&configRes)

	groupedValue["all"] = configRes.Values

	return groupedValue
}
