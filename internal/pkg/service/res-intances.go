package service

import (
	"github.com/easysoft/zendata/internal/pkg/domain"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/jinzhu/copier"
)

type ResInstancesService struct {
	ResService          *ResService          `inject:""`
	FieldService        *FieldService        `inject:""`
	CombineService      *CombineService      `inject:""`
	ResRangesService    *ResRangesService    `inject:""`
	ResInstancesService *ResInstancesService `inject:""`
}

func (s *ResInstancesService) GetResFromInstances(insts domain.ResInstances) (groupedValue map[string][]interface{}) {
	groupedValue = map[string][]interface{}{}

	for _, inst := range insts.Instances {
		for _, instField := range inst.Fields {
			s.prepareNestedInstanceRes(insts, inst, instField)
		}

		// gen values
		fieldFromInst := s.ConvertInstantsToField(insts, inst)

		s.FieldService.Generate(&fieldFromInst, true)

		group := inst.Instance
		groupedValue[group] = fieldFromInst.Values
	}

	return groupedValue
}

func (s *ResInstancesService) prepareNestedInstanceRes(insts domain.ResInstances, inst domain.ResInstancesItem, instField domain.DefField) {
	//if instField.Note == "多froms" {
	//	instField.Note += ""
	//}
	//
	//// multi froms 2.
	//if len(instField.Froms) > 0 {
	//	unionValues := make([]interface{}, 0)
	//	for _, from := range instField.Froms {
	//		if from.From == "" {
	//			from.From = instField.From
	//		}
	//
	//		from.FileDir = instField.FileDir
	//
	//		s.FieldService.Generate(&from, true)
	//
	//		unionValues = append(unionValues, from.Values...)
	//	}
	//
	//	groupedValueReferenced := map[string][]interface{}{}
	//	group := inst.Instance
	//	groupedValueReferenced[group] = unionValues
	//	vari.GlobalVars.ResData[s.ResService.GetFromKey(&instField)] = groupedValueReferenced
	//}

	// set "from" val from parent if needed
	if instField.From == "" {
		if insts.From != "" {
			instField.From = insts.From
		}
		if inst.From != "" {
			instField.From = inst.From
		}
	}

	instField.FileDir = insts.FileDir

	if instField.Use != "" { // refer to another instances or ranges
		if vari.GlobalVars.ResData[s.ResService.GetFromKey(&instField)] == nil {
			referencedRanges, referencedInstants := s.ResService.GetReferencedRangeOrInstant(instField)
			groupedValueReferenced := map[string][]interface{}{}

			if len(referencedRanges.Ranges) > 0 { // refer to ranges
				groupedValueReferenced = s.ResRangesService.GetResFromRanges(referencedRanges)

			} else if len(referencedInstants.Instances) > 0 { // refer to instances
				for _, referencedInst := range referencedInstants.Instances { // iterate records
					for _, referencedInstField := range referencedInst.Fields { // if item had children, iterate children
						s.prepareNestedInstanceRes(referencedInstants, referencedInst, referencedInstField)
					}

					field := s.ResInstancesService.ConvertInstantsToField(referencedInstants, referencedInst)

					// gen values
					group := referencedInst.Instance

					s.FieldService.Generate(&field, true)
					groupedValueReferenced[group] = field.Values
				}
			}

			vari.GlobalVars.ResData[s.ResService.GetFromKey(&instField)] = groupedValueReferenced
		}

	} else if instField.Select != "" { // refer to excel
		resFile, resType, sheet := fileUtils.GetResProp(instField.From, instField.FileDir)

		values, _ := s.ResService.GetResValueFromExcelOrYaml(resFile, resType, sheet, &instField)
		vari.GlobalVars.ResData[s.ResService.GetFromKey(&instField)] = values
	}
}

func (s *ResInstancesService) ConvertInstantsToField(insts domain.ResInstances, inst domain.ResInstancesItem) (field domain.DefField) {
	field.Note = "Converted From Instances " + insts.Title

	//field.Field = insts.Field
	field.From = insts.From

	child := domain.DefField{}
	child.Field = inst.Instance

	// some props are from parent instances
	if child.From == "" && inst.From != "" {
		child.From = inst.From
	} else if child.From == "" && insts.From != "" {
		child.From = insts.From
	}

	copier.Copy(&child, inst)

	field.Fields = append(field.Fields, child)
	field.FileDir = insts.FileDir

	return field
}
