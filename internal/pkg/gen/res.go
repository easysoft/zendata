package gen

import (
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
)

func LoadResDef(fieldsToExport []string) (res map[string]map[string][]string) {
	res = map[string]map[string][]string{}

	for index, field := range vari.Def.Fields {
		if !stringUtils.StrInArr(field.Field, fieldsToExport) {
			continue
		}

		if (field.Use != "" || field.Select != "") && field.From == "" {
			field.From = vari.Def.From
			vari.Def.Fields[index].From = vari.Def.From
		}
		loadResForFieldRecursive(&field, &res)
	}
	return
}

func loadResForFieldRecursive(field *model.DefField, res *map[string]map[string][]string) {
	if len(field.Fields) > 0 { // sub fields
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

	} else if field.From != "" && field.Type != constant.FieldTypeArticle { // from a res
		var valueMap map[string][]string
		resFile, resType, sheet := fileUtils.GetResProp(field.From, field.FileDir) // relate to current file
		valueMap, _ = getResValue(resFile, resType, sheet, field)

		if (*res)[getFromKey(field)] == nil {
			(*res)[getFromKey(field)] = map[string][]string{}
		}
		for key, val := range valueMap {
			resKey := key
			// avoid article key to be duplicate
			if vari.Def.Type == constant.ConfigTypeArticle {
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

func getResValue(resFile, resType, sheet string, field *model.DefField) (map[string][]string, string) {
	resName := ""
	groupedValues := map[string][]string{}

	if resType == "yaml" {
		groupedValues = getResFromYaml(resFile)
	} else if resType == "excel" {
		groupedValues = getResFromExcel(resFile, sheet, field)
	}

	return groupedValues, resName
}

func getResFromExcel(resFile, sheet string, field *model.DefField) map[string][]string { // , string) {
	valueMap := generateFieldValuesFromExcel(resFile, sheet, field, vari.Total)

	return valueMap
}

func getResFromYaml(resFile string) (valueMap map[string][]string) { // , resName string) {
	if vari.CacheResFileToMap[resFile] != nil { // already cached
		valueMap = vari.CacheResFileToMap[resFile]
		return
	}

	yamlContent, err := ioutil.ReadFile(resFile)
	yamlContent = stringUtils.ReplaceSpecialChars(yamlContent)

	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", resFile))
		return
	}

	insts := model.ResInstances{}
	err = yaml.Unmarshal(yamlContent, &insts)
	if err == nil && insts.Instances != nil && len(insts.Instances) > 0 { // instances
		insts.FileDir = fileUtils.GetAbsDir(resFile)
		valueMap = getResFromInstances(insts)
		//resName = insts.Field
	} else {
		ranges := model.ResRanges{}
		err = yaml.Unmarshal(yamlContent, &ranges)
		if err == nil && ranges.Ranges != nil && len(ranges.Ranges) > 0 { // ranges
			valueMap = getResFromRanges(ranges)
			//resName = ranges.Field
		} else {
			configRes := model.DefField{}
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

func getResFromInstances(insts model.ResInstances) (groupedValue map[string][]string) {
	groupedValue = map[string][]string{}

	for _, inst := range insts.Instances {
		for _, instField := range inst.Fields {
			prepareNestedInstanceRes(insts, inst, instField)
		}

		// gen values
		fieldFromInst := convertInstantToField(insts, inst)
		group := inst.Instance
		groupedValue[group] = GenerateForFieldRecursive(&fieldFromInst, false)
	}

	return groupedValue
}

func getResFromRanges(ranges model.ResRanges) map[string][]string {
	groupedValue := map[string][]string{}

	for group, expression := range ranges.Ranges {
		field := convertRangesToField(ranges, expression)

		groupedValue[group] = GenerateForFieldRecursive(&field, false)
	}

	return groupedValue
}

func prepareNestedInstanceRes(insts model.ResInstances, inst model.ResInstancesItem, instField model.DefField) {
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
				for _, referencedInst := range referencedInstants.Instances { // iterate items
					for _, referencedInstField := range referencedInst.Fields { // if item had children, iterate children
						prepareNestedInstanceRes(referencedInstants, referencedInst, referencedInstField)
					}

					field := convertInstantToField(referencedInstants, referencedInst)

					// gen values
					group := referencedInst.Instance
					groupedValueReferenced[group] = GenerateForFieldRecursive(&field, false)
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

func getReferencedRangeOrInstant(inst model.DefField) (referencedRanges model.ResRanges, referencedInsts model.ResInstances) {
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

func convertInstantToField(insts model.ResInstances, inst model.ResInstancesItem) (field model.DefField) {
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

func convertRangesToField(ranges model.ResRanges, expression string) (field model.DefField) {
	copier.Copy(&field, ranges)
	field.Range = expression

	return field
}

func getResForConfig(configRes model.DefField) map[string][]string {
	groupedValue := map[string][]string{}

	// config field is a standard field
	groupedValue["all"] = GenerateForFieldRecursive(&configRes, false)

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
