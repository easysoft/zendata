package gen

import (
	"github.com/easysoft/zendata/src/model"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func LoadResDef(fieldsToExport []string) map[string]map[string][]string {
	res := map[string]map[string][]string{}

	for _, field := range vari.Def.Fields {
		if !stringUtils.FindInArr(field.Field, fieldsToExport) { continue }

		if field.Use != "" && field.From == "" {
			field.From = vari.Def.From
		}
		loadResField(&field, &res)
	}
	return res
}

func loadResField(field *model.DefField, res *map[string]map[string][]string) {
	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			if child.Use != "" && child.From == "" {
				child.From = field.From
			}

			loadResField(&child, res)
		}
	} else if len(field.Froms) > 0 {
		for _, child := range field.Froms {
			if child.Use != "" && child.From == "" {
				child.From = field.From
			}

			loadResField(&child, res)
		}
	} else if field.From != "" {
		resFile, resType, sheet := fileUtils.GetResProp(field.From)
		values, _ := getResValue(resFile, resType, sheet, field)
		(*res)[field.From] = values

	} else if field.Config != "" {
		resFile, resType, _ := fileUtils.GetResProp(field.Config)
		values, _ := getResValue(resFile, resType, "", field)
		(*res)[field.Config] = values
	}
}

func getResValue(resFile, resType, sheet string, field *model.DefField) (map[string][]string, string) {
	resName := ""
	groupedValues := map[string][]string{}

	if resType == "yaml" {
		groupedValues, resName = getResForYaml(resFile)
	} else if resType == "excel" {
		groupedValues, resName = getResForExcel(resFile, sheet, field)
	}

	return groupedValues, resName
}

func getResForExcel(resFile, sheet string, field *model.DefField) (map[string][]string, string) {
	valueMap, resName := GenerateFieldValuesFromExcel(resFile, sheet, field)

	return valueMap, resName
}

func getResForYaml(resFile string) (valueMap map[string][]string, resName string) {
	yamlContent, err := ioutil.ReadFile(resFile)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", resFile))
		return
	}

	insts := model.ResInsts{}
	err = yaml.Unmarshal(yamlContent, &insts)
	if err == nil && insts.Instances != nil && len(insts.Instances) > 0 { // instances
		valueMap = getResForInstances(insts)
		resName = insts.Field
	} else {
		ranges := model.ResRanges{}
		err = yaml.Unmarshal(yamlContent, &ranges)
		if err == nil && ranges.Ranges != nil && len(ranges.Ranges) > 0 { // ranges
			valueMap = getResForRanges(ranges)
			resName = ranges.Field
		} else {
			configRes := model.DefField{}
			err = yaml.Unmarshal(yamlContent, &configRes)
			if err == nil { // config
				valueMap = getResForConfig(configRes)
				resName = configRes.Field
			}
		}
	}

	return
}

func getResForInstances(insts model.ResInsts) map[string][]string {
	groupedValue := map[string][]string{}

	for _, inst := range insts.Instances {
		for _, instField := range inst.Fields { // prepare referred parent instances if needed
			if instField.From == "" {
				if insts.From != "" {
					instField.From = insts.From
				}
				if inst.From != "" {
					instField.From = inst.From
				}
			}

			if instField.Use != "" { // refer to another instances or ranges
				parentRanges, parentInstants  := getRootRangeOrInstant(instField)
				groupedValueParent := map[string][]string{}

				if len(parentInstants.Instances) > 0 { // refer to instances
					for _, child := range parentInstants.Instances {
						field := convertInstantToField(parentInstants, child)

						// gen values
						group := child.Instance
						groupedValueParent[group] = GenerateForField(&field, false)
					}
				} else if len(parentRanges.Ranges) > 0 { // refer to ranges
					groupedValueParent = getResForRanges(parentRanges)
				}

				vari.Res[instField.From] = groupedValueParent
			} else if instField.Select != "" { // refer to another excel
				resFile, resType, sheet := fileUtils.GetResProp(instField.From)
				values, _ := getResValue(resFile, resType, sheet, &instField)
				vari.Res[instField.From] = values
			}
		}

		field := convertInstantToField(insts, inst)

		// gen values
		group := inst.Instance
		groupedValue[group] = GenerateForField(&field, false)
	}

	return groupedValue
}

func getRootRangeOrInstant(inst model.DefField) (parentRanges model.ResRanges, parentInsts model.ResInsts) {
	resFile, _, _ := fileUtils.GetResProp(inst.From)

	yamlContent, err := ioutil.ReadFile(resFile)
	yamlContent = ReplaceSpecialChars(yamlContent)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", resFile))
		return
	}

	err1 := yaml.Unmarshal(yamlContent, &parentRanges)
	if err1 != nil || parentRanges.Ranges == nil || len(parentRanges.Ranges) == 0 { // instances

		err2 := yaml.Unmarshal(yamlContent, &parentInsts)
		if err2 != nil || parentInsts.Instances == nil || len(parentInsts.Instances) == 0 { // ranges
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_parse_file", resFile))
			return
		}
	}

	return
}

func convertInstantToField(insts model.ResInsts, inst model.ResInst) (field model.DefField) {
	field.Field = insts.Field
	field.From = insts.From

	child := model.DefField{}
	child.Field = inst.Instance

	if child.From == "" && inst.From != "" {
		child.From = inst.From
	} else if child.From == "" && insts.From != "" {
		child.From = insts.From
	}

	copier.Copy(&child, inst)

	field.Fields = append(field.Fields, child)

	return field
}

func getResForRanges(ranges model.ResRanges) map[string][]string {
	groupedValue := map[string][]string{}

	for group, exp := range ranges.Ranges {
		// convert ranges field to standard field
		tempField := model.DefField{}
		copier.Copy(&tempField, ranges)
		tempField.Field = ranges.Field
		tempField.Range = exp

		groupedValue[group] = GenerateForField(&tempField, false)
	}

	return groupedValue
}

func getResForConfig(configRes model.DefField) map[string][]string {
	groupedValue := map[string][]string{}

	// config field is a standard field
	groupedValue["all"] = GenerateForField(&configRes, false)

	return groupedValue
}