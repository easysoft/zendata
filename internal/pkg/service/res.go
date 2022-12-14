package service

import (
	"fmt"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type ResService struct {
	FieldService *FieldService `inject:""`
	ExcelService *ExcelService `inject:""`
}

func (s *ResService) LoadResDef(fieldsToExport []string) (res map[string]map[string][]interface{}) {
	res = map[string]map[string][]interface{}{}

	for index, field := range vari.GenVars.DefData.Fields {
		if !stringUtils.StrInArr(field.Field, fieldsToExport) {
			continue
		}

		if (field.Use != "" || field.Select != "") && field.From == "" {
			field.From = vari.GenVars.DefData.From
			vari.GenVars.DefData.Fields[index].From = vari.GenVars.DefData.From
		}
		s.loadResForFieldRecursive(&field, &res)
	}
	return
}

func (s *ResService) loadResForFieldRecursive(field *model.DefField, res *map[string]map[string][]interface{}) {
	if len(field.Fields) > 0 { // child fields
		for _, child := range field.Fields {
			if child.Use != "" && child.From == "" {
				child.From = field.From
			}

			child.FileDir = field.FileDir
			s.loadResForFieldRecursive(&child, res)
		}
	} else if len(field.Froms) > 0 { // multiple from
		for _, child := range field.Froms {
			if child.Use != "" && child.From == "" {
				child.From = field.From
			}

			child.FileDir = field.FileDir
			s.loadResForFieldRecursive(&child, res)
		}

	} else if field.From != "" && field.Type != constant.FieldTypeArticle { // from a res
		var valueMap map[string][]interface{}
		resFile, resType, sheet := fileUtils.GetResProp(field.From, field.FileDir) // relate to current file
		valueMap, _ = s.getResValue(resFile, resType, sheet, field)

		if (*res)[s.getFromKey(field)] == nil {
			(*res)[s.getFromKey(field)] = map[string][]interface{}{}
		}
		for key, val := range valueMap {
			resKey := key
			// avoid article key to be duplicate
			if vari.GenVars.DefData.Type == constant.DefTypeArticle {
				resKey = resKey + "_" + field.Field
			}
			(*res)[s.getFromKey(field)][resKey] = val
		}

	} else if field.Config != "" { // from a config
		resFile, resType, _ := fileUtils.GetResProp(field.Config, field.FileDir)
		values, _ := s.getResValue(resFile, resType, "", field)
		(*res)[field.Config] = values
	}
}

func (s *ResService) getResValue(resFile, resType, sheet string, field *model.DefField) (map[string][]interface{}, string) {
	resName := ""
	groupedValues := map[string][]interface{}{}

	if resType == "yaml" {
		groupedValues = s.getResFromYaml(resFile)
	} else if resType == "excel" {
		groupedValues = s.getResFromExcel(resFile, sheet, field)
	}

	return groupedValues, resName
}

func (s *ResService) getResFromExcel(resFile, sheet string, field *model.DefField) map[string][]interface{} { // , string) {
	valueMap := s.ExcelService.generateFieldValuesFromExcel(resFile, sheet, field, vari.GenVars.Total)

	return valueMap
}

func (s *ResService) getResFromYaml(resFile string) (valueMap map[string][]interface{}) { // , resName string) {
	if vari.GenVars.CacheResFileToMap[resFile] != nil { // already cached
		valueMap = vari.GenVars.CacheResFileToMap[resFile]
		return
	}

	yamlContent, err := os.ReadFile(resFile)
	yamlContent = stringUtils.ReplaceSpecialChars(yamlContent)

	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", resFile))
		return
	}

	insts := model.ResInstances{}
	err = yaml.Unmarshal(yamlContent, &insts)
	if err == nil && insts.Instances != nil && len(insts.Instances) > 0 { // instances
		insts.FileDir = fileUtils.GetAbsDir(resFile)
		valueMap = s.getResFromInstances(insts)

	} else {
		ranges := model.ResRanges{}
		err = yaml.Unmarshal(yamlContent, &ranges)
		if err == nil && ranges.Ranges != nil && len(ranges.Ranges) > 0 { // ranges
			valueMap = s.getResFromRanges(ranges)

		} else {
			configRes := model.DefField{}
			err = yaml.Unmarshal(yamlContent, &configRes)
			if err == nil { // config
				valueMap = s.getResForConfig(configRes)

			}
		}
	}

	vari.GenVars.CacheResFileToMap[resFile] = valueMap

	return
}

func (s *ResService) getResFromInstances(insts model.ResInstances) (groupedValue map[string][]interface{}) {
	groupedValue = map[string][]interface{}{}

	//for _, inst := range insts.Instances {
	//	for _, instField := range inst.Fields {
	//		s.prepareNestedInstanceRes(insts, inst, instField)
	//	}
	//
	//	// gen values
	//	fieldFromInst := s.convertInstantToField(insts, inst)
	//	group := inst.Instance
	//	groupedValue[group] = s.GenerateForFieldRecursive(&fieldFromInst, false, vari.GenVars.Total)
	//}

	return groupedValue
}

func (s *ResService) getResFromRanges(ranges model.ResRanges) map[string][]interface{} {
	groupedValue := map[string][]interface{}{}

	//for group, expression := range ranges.Ranges {
	//	field := s.convertRangesToField(ranges, expression)
	//	groupedValue[group] = s.GenerateForFieldRecursive(&field, false, vari.GenVars.Total)
	//}

	return groupedValue
}

func (s *ResService) prepareNestedInstanceRes(insts model.ResInstances, inst model.ResInstancesItem, instField model.DefField) {
	//// set "from" val from parent if needed
	//if instField.From == "" {
	//	if insts.From != "" {
	//		instField.From = insts.From
	//	}
	//	if inst.From != "" {
	//		instField.From = inst.From
	//	}
	//}
	//instField.FileDir = insts.FileDir
	//
	//if instField.Use != "" { // refer to another instances or ranges
	//	if vari.Res[s.getFromKey(&instField)] == nil {
	//		referencedRanges, referencedInstants := s.getReferencedRangeOrInstant(instField)
	//		groupedValueReferenced := map[string][]string{}
	//
	//		if len(referencedRanges.Ranges) > 0 { // refer to ranges
	//			groupedValueReferenced = s.getResFromRanges(referencedRanges)
	//
	//		} else if len(referencedInstants.Instances) > 0 { // refer to instances
	//			for _, referencedInst := range referencedInstants.Instances { // iterate items
	//				for _, referencedInstField := range referencedInst.Fields { // if item had children, iterate children
	//					s.prepareNestedInstanceRes(referencedInstants, referencedInst, referencedInstField)
	//				}
	//
	//				field := s.convertInstantToField(referencedInstants, referencedInst)
	//
	//				// gen values
	//				group := referencedInst.Instance
	//				groupedValueReferenced[group] = s.GenerateForFieldRecursive(&field, false, vari.GenVars.Total)
	//			}
	//		}
	//
	//		vari.Res[s.getFromKey(&instField)] = groupedValueReferenced
	//	}
	//} else if instField.Select != "" { // refer to excel
	//	resFile, resType, sheet := fileUtils.GetResProp(instField.From, instField.FileDir)
	//	values, _ := s.getResValue(resFile, resType, sheet, &instField)
	//	vari.Res[s.getFromKey(&instField)] = values
	//}
}

func (s *ResService) getReferencedRangeOrInstant(inst model.DefField) (referencedRanges model.ResRanges, referencedInsts model.ResInstances) {
	resFile, _, _ := fileUtils.GetResProp(inst.From, inst.FileDir)

	yamlContent, err := ioutil.ReadFile(resFile)
	yamlContent = stringUtils.ReplaceSpecialChars(yamlContent)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", resFile))
		return
	}

	err1 := yaml.Unmarshal(yamlContent, &referencedRanges)
	if err1 != nil || referencedRanges.Ranges == nil || len(referencedRanges.Ranges) == 0 { // parse ranges failed
		err2 := yaml.Unmarshal(yamlContent, &referencedInsts)
		if err2 != nil || referencedInsts.Instances == nil || len(referencedInsts.Instances) == 0 { // parse instances failed
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_parse_file", resFile))
			return
		} else { // is instances
			referencedInsts.FileDir = fileUtils.GetAbsDir(resFile)
		}
	} else { // is ranges
		referencedRanges.FileDir = fileUtils.GetAbsDir(resFile)
	}

	return
}

func (s *ResService) convertInstantToField(insts model.ResInstances, inst model.ResInstancesItem) (field model.DefField) {
	//field.Field = insts.Field
	field.From = insts.From

	child := model.DefField{}
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

func (s *ResService) convertRangesToField(ranges model.ResRanges, expression string) (field model.DefField) {
	copier.Copy(&field, ranges)
	field.Range = expression

	return field
}

func (s *ResService) getResForConfig(configRes model.DefField) (groupedValue map[string][]interface{}) {
	groupedValue = map[string][]interface{}{}

	// config field is a standard field
	s.FieldService.GenerateValuesForNoReferField(&configRes)

	groupedValue["all"] = configRes.Values

	return groupedValue
}

func (s *ResService) getFromKey(field *model.DefField) string {
	return fmt.Sprintf("%s-%s-%s", field.From, field.Use, field.Select)
}
