package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

func LoadDefinitionFromFile(file string, fieldsToExport []string) ([]model.ClsRange, []model.ClsInst) {
	referRangeFields := make([]model.ClsRange, 0)
	referInstFields := make([]model.ClsInst, 0)

	yamlContent, err := ioutil.ReadFile(file)
	if err != nil {
		logUtils.Screen("fail to read " + file)
		return referRangeFields, referInstFields
	}

	def := model.DefData{}
	err = yaml.Unmarshal(yamlContent, &def)
	if err != nil {
		logUtils.Screen("fail to parse " + file)
		return referRangeFields, referInstFields
	}

	if strings.Index(file, "def") != 0 && constant.RootDef.Title == "" { // only add the fields in first level yaml file
		constant.RootDef = def
	}

	for _, field := range def.Fields {
		if !stringUtils.FindInArr(field.Field, fieldsToExport) { continue }

		// TODO: dealwith referRangeFields and referInstFields for constant.ResMap
	}

	return referRangeFields, referInstFields
}

func LoadReferRes([]model.ClsRange, []model.ClsInst) {
	// init const.ResMap



	//for _, field := range def.Fields {
	//	fieldValue := model.FieldValue{}
	//	// TODO: 生成fieldValue
	//	constant.ResMap[field.Field] = fieldValue // add to a map
	//}
}
