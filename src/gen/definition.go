package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func LoadRootDef(file string, fieldsToExport []string) ([]model.ClsRanges, []model.ClsInst) {
	referRangeFields := make([]model.ClsRanges, 0)
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

	constant.RootDef = def
	//constant.ResMap =

	for _, field := range def.Fields {
		if !stringUtils.FindInArr(field.Field, fieldsToExport) { continue }

		if field.From != "" {
			if field.Select != "" { // excel

			} else if field.Use != "" { // range or instance format
				//referFile, referType := getReferProp(field.From)

				// init const.ResMap
			}
		}

		// TODO:
	}

	return referRangeFields, referInstFields
}
