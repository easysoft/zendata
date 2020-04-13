package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"strings"
)

func Generate(def model.Definition, count int, fields string, out string, table string) {
	fieldArr := strings.Split(fields, ",")

	fieldMap := map[string][]interface{}{}
	for _, field := range def.Fields {
		if !stringUtils.FindInArr(field.Name, fieldArr) {
			continue
		}

		GenCol(field, count, fieldMap)
	}
}

func GenCol(field model.Field, count int, fieldMap map[string][]interface{}) {
	datatype := strings.TrimSpace(field.Datatype)
	if datatype == "" {
		datatype = "list"
	}

	switch datatype {
		case constant.LIST.String():
			GenerateList(field, count, fieldMap)

		case constant.TIMESTAMP.String():

		case constant.IP.String():

		case constant.SESSION.String():

		default:

	}
}