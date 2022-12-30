package service

import (
	"github.com/easysoft/zendata/internal/pkg/model"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

type OutputService struct {
	CombineService     *CombineService     `inject:""`
	PlaceholderService *PlaceholderService `inject:""`

	PrintService *PrintService `inject:""`
}

func (s *OutputService) GenRecords() (records []map[string]interface{}) {
	records = make([]map[string]interface{}, 0)

	for i := 0; i < vari.GlobalVars.Total; i++ {
		record := map[string]interface{}{}

		for _, field := range vari.GlobalVars.DefData.Fields {
			s.GenRecordField(&field, &record, i)
		}

		records = append(records, record)
	}

	return
}

func (s *OutputService) GenRecordField(field *model.DefField, mp *map[string]interface{}, i int) {
	if field.Join || len(field.Fields) == 0 { // set values
		val := field.Values[i%len(field.Values)]
		val = s.PlaceholderService.ReplacePlaceholder(val.(string))

		(*mp)[field.Field] = val

	} else { // set child object
		childMap := map[string]interface{}{}

		for _, child := range field.Fields {
			s.GenRecordField(&child, &childMap, i)
		}

		(*mp)[field.Field] = childMap
	}

	return
}

func (s *OutputService) PrintHumanHeaderIfNeeded() {
	if !vari.GlobalVars.Human {
		return
	}

	headerLine := ""

	for idx, field := range vari.GlobalVars.ExportFields {
		headerLine += field
		if idx < len(vari.GlobalVars.ExportFields)-1 {
			headerLine += "\t"
		}
	}

	logUtils.PrintLine(headerLine + "\n")
}
