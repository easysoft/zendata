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

func LoadClsDef(file string, fieldsToExport []string) []model.FieldValue {
	referFieldValues := make([]model.FieldValue, 0)

	yamlContent, err := ioutil.ReadFile(file)
	if err != nil {
		logUtils.Screen("fail to read " + file)
		return referFieldValues
	}

	def := model.DefData{}
	err = yaml.Unmarshal(yamlContent, &def)
	if err != nil {
		logUtils.Screen("fail to parse " + file)
		return referFieldValues
	}

	constant.RootDef = def

	for _, field := range def.Fields {
		if !stringUtils.FindInArr(field.Field, fieldsToExport) { continue }

		loadClsField(&field, &referFieldValues)
	}

	return referFieldValues
}

func loadClsField(field *model.DefField, referFieldValues *[]model.FieldValue) {
	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			loadClsField(&child, referFieldValues)
		}
	} else if field.From != "" {
		referFile, referType, tableName := getReferProp(field.From)
		fieldValue := getReferFieldValue(referFile, referType, tableName)
		*referFieldValues = append(*referFieldValues, fieldValue)
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
	} else { // excel, system.address.china
		index = strings.LastIndex(from, ".")

		left := from[:index]
		left = strings.ReplaceAll(left, ".", sep)

		referFile = left + ".xlsx"
		tableName = from[index:]
	}

	if strings.Index(referFile, "system") > -1 {
		referFile = constant.ResDir + referFile
	}

	return referFile, referType, tableName
}

func getReferFieldValue(referFile string, referType string, tableName string) model.FieldValue {
	fieldValue := model.FieldValue{}



	return fieldValue
}