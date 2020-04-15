package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"strconv"
	"strings"
)

func Generate(def *model.Definition, total int, fields string, out string, table string) [][]string {
	fieldArr := strings.Split(fields, ",")
	fieldMap := map[string][]interface{}{}
	for index, field := range def.Fields {
		if !stringUtils.FindInArr(field.Name, fieldArr) {
			continue
		}
		GenerateFieldArr(&field, total, fieldMap)
		def.Fields[index].Precision = field.Precision
	}

	rows := make([][]string, 0)
	for i := 0; i < total; i++ {
		for _, field := range def.Fields {
			if !stringUtils.FindInArr(field.Name, fieldArr) {
				continue
			}

			if len(rows) == i { rows = append(rows, make([]string, 0)) }

			str := "n/a"
			val := fieldMap[field.Name][i]
			switch val.(type) {
				case int64:
					str = strconv.FormatInt(val.(int64),10)
				case float64:
					precision := -1
					if field.Precision >= 0 {
						precision = field.Precision
					}
					str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
				case byte:
					str = string(val.(byte))
				default:
			}

			rows[i] = append(rows[i], str)
		}
	}

	return rows
}

func GenerateFieldArr(field *model.Field, total int, fieldMap map[string][]interface{}) {
	datatype := strings.TrimSpace(field.Datatype)
	if datatype == "" {
		datatype = "list"
	}

	switch datatype {
		case constant.LIST.String():
			GenerateList(field, total, fieldMap)

		case constant.TIMESTAMP.String():
			GenerateTimestamp(field, total, fieldMap)

		case constant.IP.String():
			GenerateIP(field, total, fieldMap)

		case constant.SESSION_ID.String():
			GenerateSessionId(field, total, fieldMap)

		default:
	}
}