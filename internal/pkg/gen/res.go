package gen

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/helper"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func LoadResDef(fieldsToExport []string) (res map[string]map[string][]string) {
	res = map[string]map[string][]string{}

	for index, field := range vari.GlobalVars.DefData.Fields {
		if !stringUtils.StrInArr(field.Field, fieldsToExport) {
			continue
		}

		if (field.Use != "" || field.Select != "") && field.From == "" {
			field.From = vari.GlobalVars.DefData.From
			vari.GlobalVars.DefData.Fields[index].From = vari.GlobalVars.DefData.From
		}
		loadResForFieldRecursive(&field, &res)
	}
	return
}

func loadResForFieldRecursive(field *domain.DefField, res *map[string]map[string][]string) {
	if len(field.Fields) > 0 { // child fields
		for _, child := range field.Fields {
			if child.Use != "" && child.From == "" {
				child.From = field.From
			}

			child.FileDir = field.FileDir
			loadResForFieldRecursive(&child, res)
		}
	} else if len(field.Froms) > 0 { // multiple from
		for _, child := range field.Froms {
			if child.Use != "" && child.From == "" {
				child.From = field.From
			}

			child.FileDir = field.FileDir
			loadResForFieldRecursive(&child, res)
		}

	} else if field.From != "" && field.Type != consts.FieldTypeArticle { // from a res
		var valueMap map[string][]string
		resFile, resType, sheet := fileUtils.GetResProp(field.From, field.FileDir) // relate to current file
		valueMap, _ = getResValue(resFile, resType, sheet, field)

		if (*res)[getFromKey(field)] == nil {
			(*res)[getFromKey(field)] = map[string][]string{}
		}
		for key, val := range valueMap {
			resKey := key
			// avoid article key to be duplicate
			if vari.GlobalVars.DefData.Type == consts.DefTypeArticle {
				resKey = resKey + "_" + field.Field
			}
			(*res)[getFromKey(field)][resKey] = val
		}

	} else if field.Config != "" { // from a config
		resFile, resType, _ := fileUtils.GetResProp(field.Config, field.FileDir)
		values, _ := getResValue(resFile, resType, "", field)
		(*res)[field.Config] = values
	}
}

func getResValue(resFile, resType, sheet string, field *domain.DefField) (map[string][]string, string) {
	resName := ""
	groupedValues := map[string][]string{}

	if resType == "yaml" {
		groupedValues = getResFromYaml(resFile)
	} else if resType == "excel" {
		groupedValues = getResFromExcel(resFile, sheet, field)
	}

	return groupedValues, resName
}

func getResFromExcel(resFile, sheet string, field *domain.DefField) map[string][]string { // , string) {
	valueMap := generateFieldValuesFromExcel(resFile, sheet, field, vari.GlobalVars.Total)

	return valueMap
}

func getResFromYaml(resFile string) (valueMap map[string][]string) { // , resName string) {
	if vari.CacheResFileToMap[resFile] != nil { // already cached
		valueMap = vari.CacheResFileToMap[resFile]
		return
	}

	yamlContent, err := ioutil.ReadFile(resFile)
	yamlContent = helper.ReplaceSpecialChars(yamlContent)

	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", resFile))
		return
	}

	insts := domain.ResInstances{}
	err = yaml.Unmarshal(yamlContent, &insts)
	if err == nil && insts.Instances != nil && len(insts.Instances) > 0 { // instances
		insts.FileDir = fileUtils.GetAbsDir(resFile)
		valueMap = getResFromInstances(insts)
		//resName = insts.Field

	} else {
		ranges := domain.ResRanges{}
		err = yaml.Unmarshal(yamlContent, &ranges)
		if err == nil && ranges.Ranges != nil && len(ranges.Ranges) > 0 { // ranges
			valueMap = getResFromRanges(ranges)
			//resName = ranges.Field

		} else {
			configRes := domain.DefField{}
			err = yaml.Unmarshal(yamlContent, &configRes)
			if err == nil { // config
				valueMap = getResForConfig(configRes)
				//resName = configRes.Field

			}
		}
	}

	vari.CacheResFileToMap[resFile] = valueMap
	//vari.CacheResFileToName[resFile] = resName

	return
}

func getResFromInstances(insts domain.ResInstances) (groupedValue map[string][]string) {
	groupedValue = map[string][]string{}

	for _, inst := range insts.Instances {
		for _, instField := range inst.Fields {
			prepareNestedInstanceRes(insts, inst, instField)
		}

		// gen values
		fieldFromInst := convertInstantToField(insts, inst)
		group := inst.Instance
		groupedValue[group] = GenerateForFieldRecursive(&fieldFromInst, false, vari.GlobalVars.Total)
	}

	return groupedValue
}

func getResFromRanges(ranges domain.ResRanges) map[string][]string {
	groupedValue := map[string][]string{}

	for group, expression := range ranges.Ranges {
		field := convertRangesToField(ranges, expression)
		groupedValue[group] = GenerateForFieldRecursive(&field, false, vari.GlobalVars.Total)
	}

	return groupedValue
}

func prepareNestedInstanceRes(insts domain.ResInstances, inst domain.ResInstancesItem, instField domain.DefField) {
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
		if vari.Res[getFromKey(&instField)] == nil {
			referencedRanges, referencedInstants := getReferencedRangeOrInstant(instField)
			groupedValueReferenced := map[string][]string{}

			if len(referencedRanges.Ranges) > 0 { // refer to ranges
				groupedValueReferenced = getResFromRanges(referencedRanges)

			} else if len(referencedInstants.Instances) > 0 { // refer to instances
				for _, referencedInst := range referencedInstants.Instances { // iterate records
					for _, referencedInstField := range referencedInst.Fields { // if item had children, iterate children
						prepareNestedInstanceRes(referencedInstants, referencedInst, referencedInstField)
					}

					field := convertInstantToField(referencedInstants, referencedInst)

					// gen values
					group := referencedInst.Instance
					groupedValueReferenced[group] = GenerateForFieldRecursive(&field, false, vari.GlobalVars.Total)
				}
			}

			vari.Res[getFromKey(&instField)] = groupedValueReferenced
		}
	} else if instField.Select != "" { // refer to excel
		resFile, resType, sheet := fileUtils.GetResProp(instField.From, instField.FileDir)
		values, _ := getResValue(resFile, resType, sheet, &instField)
		vari.Res[getFromKey(&instField)] = values
	}
}

func getReferencedRangeOrInstant(inst domain.DefField) (referencedRanges domain.ResRanges, referencedInsts domain.ResInstances) {
	resFile, _, _ := fileUtils.GetResProp(inst.From, inst.FileDir)

	yamlContent, err := ioutil.ReadFile(resFile)
	yamlContent = helper.ReplaceSpecialChars(yamlContent)
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

func convertInstantToField(insts domain.ResInstances, inst domain.ResInstancesItem) (field domain.DefField) {
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

func convertRangesToField(ranges domain.ResRanges, expression string) (field domain.DefField) {
	copier.Copy(&field, ranges)
	field.Range = expression

	return field
}

func getResForConfig(configRes domain.DefField) map[string][]string {
	groupedValue := map[string][]string{}

	// config field is a standard field
	groupedValue["all"] = GenerateForFieldRecursive(&configRes, false, vari.GlobalVars.Total)

	return groupedValue
}

//func getLastDuplicateVal(preMap map[string][]string, key string) (valMap map[string][]string) {
//	lastKey := ""
//	for k := range preMap {
//		if key == removeKeyNumber(k) {
//			lastKey = k
//			break
//		}
//	}
//
//	if lastKey == "" || preMap[lastKey] == nil {
//		return nil
//	}
//
//	valMap = map[string][]string{}
//	valMap[key] = preMap[lastKey]
//	return
//}
//func removeKeyNumber(key string) string {
//	arr := strings.Split(key, "_")
//	ret := strings.Join(arr[:len(arr)-1], "_")
//	return ret
//}
