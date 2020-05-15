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

func LoadRootDef(file string, fieldsToExport []string) ([]model.ClsRange, []model.ClsInst) {
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

	constant.RootDef = def
	//constant.ResMap =

	for _, field := range def.Fields {
		if !stringUtils.FindInArr(field.Field, fieldsToExport) { continue }

		if field.From != "" {
			if field.Select != "" { // excel

			} else if field.Use != "" { // range or instance format
				//referFile, referType := getReferPath(field.From)

				// init const.ResMap
			}
		}

		// TODO:
	}

	return referRangeFields, referInstFields
}

func getReferPath(from string) (string, string, string) {
	referFile := ""
	referType := ""
	tableName := ""

	sep := string(os.PathSeparator)

	index := strings.LastIndex(from, ".yaml")
	if index > -1 { // system.nubmer.yaml
		left := from[:index]
		left = strings.ReplaceAll(left, ".", sep)

		referFile = left + ".yaml"
	} else { // system.address.china
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