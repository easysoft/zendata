package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"strings"
)

func Generate(def model.Definition, count int, fields string, out string, table string) {
	content := ""
	for i := 0; i < count; i++ {
		AddRow(def, fields, &content)
	}
}

func AddRow(def model.Definition, fields string, content *string) {
	fieldArr := strings.Split(fields, ",")

	for _, field := range def.Fields {
		if !stringUtils.FindInArr(field.Name, fieldArr) {
			continue
		}

		AddCol(field, content)
		*content = *content + "\r\n"
	}
}

func AddCol(field model.Field, content *string) {
	datatype := strings.TrimSpace(field.Datatype)
	if datatype == "" {
		datatype = "list"
	}

	switch datatype {
		case constant.LIST.String():
			GenerateList(field, content)

		case constant.TIMESTAMP.String():

		case constant.IP.String():

		case constant.SESSION.String():

		default:

	}
}