package service

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
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

func (s *OutputService) GenRecordField(field *domain.DefField, mp *map[string]interface{}, i int) {
	if field.Join || len(field.Fields) == 0 { // set values
		val := field.Values[i%len(field.Values)]

		switch val.(type) {
		case string:
			val = s.PlaceholderService.ReplacePlaceholder(fmt.Sprintf("%v", val))
		default:
		}

		(*mp)[field.Field] = val

	} else { // set children
		var childVal interface{}

		isRecursive := field.Mode == consts.ModeRecursive || field.Mode == consts.ModeRecursiveShort
		indexArr := make([]int, 0)
		if isRecursive {
			indexArr = s.getModArrForChildrenRecursive(field)
		}

		if field.Items == 0 { // output is object
			mp := map[string]interface{}{}
			for k, child := range field.Fields {
				var index int
				if isRecursive {
					mod := indexArr[k]
					index = i / mod % len(child.Values)
				} else {
					index = i % len(child.Values)
				}

				s.GenRecordField(&child, &mp, index)
			}

			childVal = mp

		} else { // output is array
			var mpArr []map[string]interface{}

			for itemIndex := 0; itemIndex < field.Items; itemIndex++ {
				mp := map[string]interface{}{}
				for k, child := range field.Fields {
					index := i*field.Items + itemIndex

					if isRecursive {
						mod := indexArr[k]
						index = index / mod % len(child.Values)
					} else {
						index = index % len(child.Values)
					}

					s.GenRecordField(&child, &mp, index)
				}

				mpArr = append(mpArr, mp)
			}

			childVal = mpArr
		}

		(*mp)[field.Field] = childVal
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

func (s *OutputService) getModArrForChildrenRecursive(field *domain.DefField) []int {
	indexArr := make([]int, 0)
	for _, _ = range field.Fields {
		indexArr = append(indexArr, 0)
	}

	for i := 0; i < len(field.Fields); i++ {
		loop := 1
		for j := i + 1; j < len(field.Fields); j++ {
			loop = loop * len(field.Fields[j].Values)

			//if field.Items > 1 {
			//	loop /= field.Items
			//}
		}

		indexArr[i] = loop
	}

	return indexArr
}
