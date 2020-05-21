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
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func LoadResDef(fieldsToExport []string) map[string]map[string][]string {
	res := map[string]map[string][]string{}

	for _, field := range constant.Def.Fields {
		if !stringUtils.FindInArr(field.Field, fieldsToExport) { continue }

		loadResField(&field, &res)
	}

	return res
}

func loadResField(field *model.DefField, res *map[string]map[string][]string) {
	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			loadResField(&child, res)
		}
	} else if field.From != "" {
		resFile, resType := getResProp(field.From)
		values, _ := getResValue(resFile, resType, field)
		(*res)[field.From] = values
	}
}

func getResProp(from string) (string, string) {
	resFile := ""
	resType := ""

	sep := string(os.PathSeparator)

	index := strings.LastIndex(from, ".yaml")
	if index > -1 { // yaml, system.ip.v1.yaml
		left := from[:index]
		left = strings.ReplaceAll(left, ".", sep)

		resFile = left + ".yaml"
		resType = "yaml"
	} else { // excel, system.address.v1.city
		index = strings.LastIndex(from, ".")

		left := from[:index]
		left = strings.ReplaceAll(left, ".", sep)

		resFile = left + ".xlsx"
		resType = "excel"
	}

	if strings.Index(resFile, "system") == -1 { // no system cls
		resFile = vari.InputDir + resFile

		if !fileUtils.FileExist(resFile) { // not in input dir
			resFile = vari.ExeDir + resFile

			if !fileUtils.FileExist(resFile) { // not in exe dir
				resFile = ""
			}
		}
	} else {
		resFile = constant.ResDir + resFile
	}

	return resFile, resType
}

func getResValue(resFile string, resType string, field *model.DefField) (map[string][]string, string) {
	resName := ""
	groupedValues := map[string][]string{}

	if resType == "yaml" {
		groupedValues, resName = getResForYaml(resFile)
	} else if resType == "excel" {
		groupedValues, resName = getResForExcel(resFile, field)
	}

	return groupedValues, resName
}

func getResForExcel(resFile string, field *model.DefField) (map[string][]string, string) {
	valueMap, resName := GenerateFieldValuesFromExcel(resFile, field)

	return valueMap, resName
}

func getResForYaml(resFile string) (map[string][]string, string) {
	resName := ""
	valueMap := map[string][]string{}

	ranges := model.ResRanges{}

	yamlContent, err := ioutil.ReadFile(resFile)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", resFile))
		return valueMap, ""
	}

	err = yaml.Unmarshal(yamlContent, &ranges)
	if err != nil || ranges.Ranges == nil || len(ranges.Ranges) == 0 { // instances
		if vari.Verbose { logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_parse_ranges", resFile)) }

		insts := model.ResInsts{}
		err = yaml.Unmarshal(yamlContent, &insts)
		if err != nil {
			return valueMap, ""
		} else {
			valueMap = getResForInstances(insts)
		}

		resName = insts.Field

	} else { // ranges
		valueMap = getResForRanges(ranges)
		resName = ranges.Field
	}

	return valueMap, resName
}

func getResForRanges(ranges model.ResRanges) map[string][]string {
	groupedValue := map[string][]string{}

	for group, exp := range ranges.Ranges {
		// convert ranges field to standard field
		tempField := model.DefField{}
		copier.Copy(&tempField, ranges)
		tempField.Field = ranges.Field
		tempField.Range = exp
		tempField.Type = "cls"

		groupedValue[group] = GenerateForField(&tempField, constant.Total)
	}

	return groupedValue
}

func getResForInstances(insts model.ResInsts) map[string][]string {
	groupedValue := map[string][]string{}

	for _, inst := range insts.Instances {
		group := inst.Instance

		// convert instant field to standard field
		tempField := model.DefField{Field: insts.Field, Type: "cls"}

		child := model.DefField{Field: inst.Instance}
		copier.Copy(&child, inst)

		tempField.Fields = append(tempField.Fields, child)

		groupedValue[group] = GenerateForField(&tempField, constant.Total)
	}

	return groupedValue
}