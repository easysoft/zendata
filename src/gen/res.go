package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
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
		resFile, resType, sheet := getResProp(field.From)
		values, _ := getResValue(resFile, resType, sheet, field)
		(*res)[field.From] = values

	} else if field.Config != "" {
		resFile, resType, _ := getResProp(field.Config)
		values, _ := getResValue(resFile, resType, "", field)
		(*res)[field.Config] = values
	}
}

func getResProp(from string) (resFile, resType, sheet string) { // from resource

	index := strings.LastIndex(from, ".yaml")
	if index > -1 { // yaml, ip.v1.yaml
		resFile = convertYamlPath(from)
		resType = "yaml"
	} else { // excel, like address.cn.v1.china
		resFile, sheet = convertExcelPath(from)
		resType = "excel"
	}

	if resFile == "" {
		resPath := vari.ConfigDir + resFile
		if !fileUtils.FileExist(resPath) { // in same folder with passed config file

			resPath = vari.WorkDir + resFile
			if !fileUtils.FileExist(resPath) {  // in res file
				resPath = ""
			}
		}
		resFile = resPath
	}

	return
}

func convertYamlPath(from string) (ret string) {
	arr := strings.Split(from, ".")
	for i := 0; i < len(arr); i++ {
		dir := ""
		if i > 0 {
			dir = strings.Join(arr[:i], constant.PthSep)
		}
		file := strings.Join(arr[i:], ".")

		if dir != "" {
			ret = dir + constant.PthSep + file
		} else {
			ret = file
		}

		realPth1 := vari.WorkDir + constant.ResDirYaml + ret
		realPth2 := vari.WorkDir + constant.ResDirUsers + ret
		if fileUtils.FileExist(realPth1) {
			ret = realPth1
			break
		} else if fileUtils.FileExist(realPth2) {
			ret = realPth2
			break
		}
	}

	return
}

func convertExcelPath(from string) (ret, sheet string) {
	path1 := from // address.cn.v1
	index := strings.LastIndex(from, ".")
	path2 := from[:index] // address.cn.v1.china

	paths := [2]string{path1, path2}
	for index, path := range paths {

		arr := strings.Split(path, ".")
		for i := 0; i < len(arr); i++ {
			dir := ""
			if i > 0 {
				dir = strings.Join(arr[:i], constant.PthSep)
			}
			file := strings.Join(arr[i:], ".") + ".xlsx"

			if dir != "" {
				ret = dir + constant.PthSep + file
			} else {
				ret = file
			}

			realPth := vari.WorkDir + constant.ResDirData + ret
			if fileUtils.FileExist(realPth) {
				if index == 1 {
					sheet = from[strings.LastIndex(from, ".")+1:]
				}
				ret = realPth
				return
			}
		}
	}

	return
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
			if instField.Use != "" { // refer to another def
				if instField.From == "" {
					if insts.From != "" {
						instField.From = insts.From
					}
					if inst.From != "" {
						instField.From = inst.From
					}
				}

				parentRanges, parentInstants  := getRootRangeOrInstant(instField)
				groupedValueParent := map[string][]string{}

				if len(parentInstants.Instances) > 0 {
					for _, child := range parentInstants.Instances {
						field := convertInstantToField(parentInstants, child)

						// gen values
						group := child.Instance
						groupedValueParent[group] = GenerateForField(&field, vari.Total, false)
					}
				} else if len(parentRanges.Ranges) > 0 {
					groupedValueParent = getResForRanges(parentRanges)
				}

				vari.Res[instField.From] = groupedValueParent
			}
		}

		field := convertInstantToField(insts, inst)

		// gen values
		group := inst.Instance
		groupedValue[group] = GenerateForField(&field, vari.Total, false)
	}

	return groupedValue
}

func getRootRangeOrInstant(inst model.DefField) (parentRanges model.ResRanges, parentInsts model.ResInsts) {
	resFile, _, _ := getResProp(inst.From)

	yamlContent, err := ioutil.ReadFile(resFile)
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

		groupedValue[group] = GenerateForField(&tempField, vari.Total, false)
	}

	return groupedValue
}

func getResForConfig(configRes model.DefField) map[string][]string {
	groupedValue := map[string][]string{}

	// config field is a standard field
	groupedValue["all"] = GenerateForField(&configRes, vari.Total, false)

	return groupedValue
}