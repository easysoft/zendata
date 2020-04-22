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
	indexOfRow := 0
	for i := 0; i < total; i++ {
		for _, field := range def.Fields {
			if !stringUtils.FindInArr(field.Name, fieldArr) {
				continue
			}

			if field.Loop == 0 {
				field.Loop = 1
			}

			prefix := field.Prefix
			postfix := field.Postfix
			if len(rows) == i { rows = append(rows, make([]string, 0)) }

			loopStr := ""
			for j := 0; j < field.Loop; j++ {
				if loopStr != "" {
					loopStr = loopStr + field.Loopfix
				}
				str := GetFieldStr(field, fieldMap, indexOfRow)
				loopStr = loopStr + str

				indexOfRow++
				if indexOfRow == len(fieldMap[field.Name]) { // no enough
					indexOfRow = 0
				}
			}

			rows[i] = append(rows[i], prefix + loopStr + postfix)
		}
	}

	return rows
}

func GenerateFieldArr(field *model.Field, total int, fieldMap map[string][]interface{}) {
	datatype := strings.TrimSpace(field.Type)
	if datatype == "" { datatype = "list" }

	fieldValue := model.FieldValue{}

	switch datatype {
	case constant.LIST.String():
		fieldValue = GenerateList(field, total)

	default:
	}

	fieldMap[field.Name] = GetFieldPlatArr(fieldValue, total)
}

func GetFieldPlatArr(fieldValue model.FieldValue, total int) []interface{} {
	arr := make([]interface{}, 0)

	if len(fieldValue.Children) > 0 {

	} else {
		arr = append(arr, fieldValue.Values...)
	}

	return arr
}

func GetFieldStr(field model.Field, fieldMap map[string][]interface{}, indexOfRow int) string {
	str := "n/a"
	success := false
	val := fieldMap[field.Name][indexOfRow]
	switch val.(type) {
		case int64:
			if field.Format != "" {
				str, success = stringUtils.FormatStr(field.Format, val.(int64))
			}
			if !success {
				str = strconv.FormatInt(val.(int64), 10)
			}
		case float64:
			precision := 0
			if field.Precision > 0 {
				precision = field.Precision
			}
			if field.Format != "" {
				str, success = stringUtils.FormatStr(field.Format, val.(float64))
			}
			if !success {
				str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
			}
		case byte:
			str = string(val.(byte))
			if field.Format != "" {
				str, success = stringUtils.FormatStr(field.Format, str)
			}
			if !success {
				str = string(val.(byte))
			}
		default:
	}

	return str
}