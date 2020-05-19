package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func LoadClsDef(file string, fieldsToExport []string) map[string]map[string][]string {
	referFieldValueMap := map[string]map[string][]string{}

	for _, field := range constant.RootDef.Fields {
		if !stringUtils.FindInArr(field.Field, fieldsToExport) { continue }

		loadClsField(&field, &referFieldValueMap)
	}

	return referFieldValueMap
}

func loadClsField(field *model.DefField, referFieldValueMap *map[string]map[string][]string) {
	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			loadClsField(&child, referFieldValueMap)
		}
	} else if field.From != "" {
		referFile, referType := getReferProp(field.From)
		values, _ := getReferFieldValue(referFile, referType, field)
		(*referFieldValueMap)[field.From] = values
	}
}

func getReferProp(from string) (string, string) {
	referFile := ""
	referType := ""

	sep := string(os.PathSeparator)

	index := strings.LastIndex(from, ".yaml")
	if index > -1 { // yaml, system.nubmer.yaml
		left := from[:index]
		left = strings.ReplaceAll(left, ".", sep)

		referFile = left + ".yaml"
		referType = "yaml"
	} else { // excel, system.address.china
		index = strings.LastIndex(from, ".")

		left := from[:index]
		left = strings.ReplaceAll(left, ".", sep)

		referFile = left + ".xlsx"
		referType = "excel"
	}

	if strings.Index(referFile, "system") == -1 { // no system cls
		referFile = vari.InputDir + referFile

		if !fileUtils.FileExist(referFile) { // not in input dir
			referFile = vari.ExeDir + referFile

			if !fileUtils.FileExist(referFile) { // not in exe dir
				referFile = ""
			}
		}
	} else {
		referFile = constant.ResDir + referFile
	}

	return referFile, referType
}

func getReferFieldValue(referFile string, referType string, field *model.DefField) (map[string][]string, string) {
	name := ""
	values := map[string][]string{}

	if referType == "yaml" {
		values, name = getReferFieldValueForYaml(referFile)
	} else if referType == "excel" {
		values, name = getReferFieldValueForExcel(referFile, field)
	}

	return values, name
}

func getReferFieldValueForYaml(referFile string) (map[string][]string, string) {
	name := ""
	valueMap := map[string][]string{}

	ranges := model.ClsRanges{}

	yamlContent, err := ioutil.ReadFile(referFile)
	if err != nil {
		logUtils.Screen("fail to read " + referFile)
		return valueMap, ""
	}

	err = yaml.Unmarshal(yamlContent, &ranges)
	if err != nil || ranges.Ranges == nil || len(ranges.Ranges) == 0 { // instances
		logUtils.Screen("not ClsRanges " + referFile + ", try to parse as ClsInsts")

		insts := model.ClsInsts{}
		err = yaml.Unmarshal(yamlContent, &insts)
		if err != nil {
			return valueMap, ""
		} else {
			valueMap = getReferFieldValueForInstances(insts)
		}

		name = insts.Field

	} else { // ranges
		valueMap = getReferFieldValueForRanges(ranges)
		name = ranges.Field
	}

	return valueMap, name
}

func getReferFieldValueForRanges(ranges model.ClsRanges) map[string][]string {
	values := map[string][]string{}

	for name, exp := range ranges.Ranges {
		// convert ranges refer to standard field
		tempField := model.DefField{Field: ranges.Field, Range: exp, Type: "cls"}

		values[name] = GenerateFieldItemsFromDefinition(&tempField)
	}

	return values
}

func getReferFieldValueForInstances(insts model.ClsInsts) map[string][]string {
	values := map[string][]string{}

	for _, inst := range insts.Instances {
		// convert ranges refer to standard field
		tempField := model.DefField{Field: insts.Field, Type: "cls"}

		child := model.DefField{Field: inst.Instance}
		copier.Copy(&child, inst)

		tempField.Fields = append(tempField.Fields, child)

		values[inst.Instance] = GenerateForField(&tempField, constant.Total)
	}

	return values
}

func getReferFieldValueForExcel(referFile string, field *model.DefField) (map[string][]string, string) {
	values, dbName := GenerateFieldValuesFromExcel(referFile, field)

	return values, dbName
}