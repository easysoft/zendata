package service

import (
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

type ResExcelService struct {
	ExcelService *ExcelService `inject:""`
}

func (s *ResExcelService) GetResFromExcel(resFile, sheet string, field *domain.DefField) map[string][]interface{} { // , string) {
	valueMap := s.ExcelService.generateFieldValuesFromExcel(resFile, sheet, field, vari.GlobalVars.Total)

	return valueMap
}
