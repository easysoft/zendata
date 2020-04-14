package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"strconv"
	"strings"
)

func Generate(def model.Definition, total int, fields string, out string, table string) [][]string {
	fieldArr := strings.Split(fields, ",")
	fieldMap := map[string][]interface{}{}
	for _, field := range def.Fields {
		if !stringUtils.FindInArr(field.Name, fieldArr) {
			continue
		}
		GenerateFieldArr(field, total, fieldMap)
	}

	rows := make([][]string, 0)
	for i := 0; i < total; i++ {
		for _, field := range def.Fields {
			if !stringUtils.FindInArr(field.Name, fieldArr) {
				continue
			}

			str := fieldMap[field.Name][i].(int64)

			if len(rows) == i {
				rows = append(rows, make([]string, 0))
			}
			rows[i] = append(rows[i], strconv.FormatInt(str,10))
		}
	}

	return rows
}

func GenerateFieldArr(field model.Field, total int, fieldMap map[string][]interface{}) {
	datatype := strings.TrimSpace(field.Datatype)
	if datatype == "" {
		datatype = "list"
	}

	switch datatype {
		case constant.LIST.String():
			GenerateList(field, total, fieldMap)

		case constant.TIMESTAMP.String():

		case constant.IP.String():

		case constant.SESSION.String():

		default:

	}
}