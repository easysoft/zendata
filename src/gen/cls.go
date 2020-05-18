package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func LoadClsDef(file string, fieldsToExport []string) map[string]map[string][]string {
	referFieldValueMap := map[string]map[string][]string{}

	yamlContent, err := ioutil.ReadFile(file)
	if err != nil {
		logUtils.Screen("fail to read " + file)
		return referFieldValueMap
	}

	def := model.DefData{}
	err = yaml.Unmarshal(yamlContent, &def)
	if err != nil {
		logUtils.Screen("fail to parse " + file)
		return referFieldValueMap
	}

	constant.RootDef = def

	for _, field := range def.Fields {
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
		referFile, referType, tableName := getReferProp(field.From)
		values := getReferFieldValue(referFile, referType, tableName)
		(*referFieldValueMap)[field.Field] = values
	}
}

func getReferProp(from string) (string, string, string) {
	referFile := ""
	referType := ""
	tableName := ""

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
		tableName = from[index:]
		referType = "excel"
	}

	if strings.Index(referFile, "system") > -1 {
		referFile = constant.ResDir + referFile
	}

	return referFile, referType, tableName
}

func getReferFieldValue(referFile string, referType string, tableName string) map[string][]string {
	values := map[string][]string{}

	if referType == "yaml" {
		values = getReferFieldValueForYaml(referFile)
	} else if referType == "excel" {
		values = getReferFieldValueForExcel(referFile, tableName)
	}

	return values
}

func getReferFieldValueForYaml(referFile string) map[string][]string {
	valueMap := map[string][]string{}

	ranges := model.ClsRanges{}

	yamlContent, err := ioutil.ReadFile(referFile)
	if err != nil {
		logUtils.Screen("fail to read " + referFile)
		return valueMap
	}

	err = yaml.Unmarshal(yamlContent, &ranges)
	if err != nil || ranges.Ranges == nil || len(ranges.Ranges) == 0 {
		logUtils.Screen("fail to parse ClsRanges " + referFile + ", try to parse as ClsInsts")

		insts := model.ClsInsts{}

		yamlContent, err := ioutil.ReadFile(referFile)
		if err != nil {
			logUtils.Screen("fail to read " + referFile)
			return valueMap
		}

		err = yaml.Unmarshal(yamlContent, &insts)
		if err != nil {
			return valueMap
		} else {
			valueMap = getReferFieldValueForInstances(insts)
		}

	} else {
		valueMap = getReferFieldValueForRanges(ranges)
	}

	return valueMap
}

func getReferFieldValueForRanges(ranges model.ClsRanges) map[string][]string {
	values := map[string][]string{}

	for name, exp := range ranges.Ranges {
		// convert ranges refer to standard field
		tempField := model.DefField{Field: ranges.Field, Range: exp}

		values[name] = GenerateFieldItemsFromDefinition(&tempField, constant.Total)
	}

	return values
}

func getReferFieldValueForInstances(referFile model.ClsInsts) map[string][]string {
	values := map[string][]string{}



	return values
}

func getReferFieldValueForExcel(referFile string, tableName string) map[string][]string {
	values := map[string][]string{}



	return values
}