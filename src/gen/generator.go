package gen

import (
	"encoding/json"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
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

			str := GenerateFieldValWithLoop(field, fieldMap, indexOfRow)
			if len(rows) == i { rows = append(rows, make([]string, 0)) }
			rows[i] = append(rows[i], str)
			indexOfRow++
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

	fieldMap[field.Name] = append(fieldMap[field.Name], fieldValue)
}

func GenerateFieldValWithLoop(field model.Field, fieldMap map[string][]interface{}, indexOfRow int) string {
	prefix := field.Prefix
	postfix := field.Postfix

	loopStr := ""
	for j := 0; j < field.Loop; j++ {
		if loopStr != "" {
			loopStr = loopStr + field.Loopfix
		}
		str := GenerateFieldVal(field, fieldMap, indexOfRow)
		loopStr = loopStr + str

		indexOfRow++
	}

	return prefix + loopStr + postfix
}

func GenerateFieldVal(field model.Field, fieldMap map[string][]interface{}, indexOfRow int) string {
	fieldValue := fieldMap[field.Name][indexOfRow % len(fieldMap[field.Name])].(model.FieldValue)

	str := ""
	index := indexOfRow % len(fieldValue.Values)
	if len(fieldValue.Children) == 0 {
		val := fieldValue.Values[index]
		str = GetFieldValStr(field, val)
	} else {
		str = GetFieldValStrFromNestedObj(fieldValue, index)
	}

	return str
}

func GetFieldValStrFromNestedObj(fieldValue model.FieldValue, index int) string {
	arr := GetFieldPlatArr(fieldValue)

	bytes, _ := json.Marshal(arr)
	logUtils.Screen(string(bytes))

	return ""
}

func GetFieldPlatArr(fieldValue model.FieldValue) []interface{} {
	arr := make([]interface{}, 0)

	if len(fieldValue.Children) > 0 {
		platArr := ConvertNestedFieldToPlatArr(fieldValue)
		arr = append(arr, platArr...)
	} else {
		arr = append(arr, fieldValue.Values...)
	}

	return arr
}

func ConvertNestedFieldToPlatArr(fieldValue model.FieldValue) []interface{} {
	arr := make([]interface{}, 0)

	for _, child := range fieldValue.Children {
		if child.Level == 0 {

		}
	}

	return arr
}

func GetFieldValStr(field model.Field, val interface{}) string {
	str := "n/a"
	success := false

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

