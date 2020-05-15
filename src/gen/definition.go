package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
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

	if strings.Index(file, "def") != 0 && constant.Definition.Title == "" { // only add the fields in first level yaml file
		constant.Definition = def
	}

	for _, field := range def.Fields {
		_ = field
		// TODO: dealwith referRangeFields and referInstFields for constant.LoadedResValues
	}

	//for _, field := range def.Fields {
	//	fieldValue := model.FieldValue{}
	//	// TODO: 生成fieldValue
	//	constant.LoadedResValues[field.Field] = fieldValue // add to a map
	//}

	return referRangeFields, referInstFields
}

func LoadReferRes([]model.ClsRange, []model.ClsInst) {
	//fieldValueMap := map[string]model.FieldValue{}
	//
	//LoadDefinitionFromFile(constant.ResBuildIn) // load buildin resource
}
