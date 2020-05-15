package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"strconv"
	"strings"
)

func GenerateForDefinition(total int, fieldsToExport []string) ([][]string, []bool) {
	def := constant.Definition

	fieldNameToValues := map[string][]string{}

	colTypes := make([]bool, 0)

	// 为每个field生成值列表
	for index, field := range def.Fields {
		if !stringUtils.FindInArr(field.Field, fieldsToExport) {
			continue
		}

		values := GenerateForField(&field, total)
		def.Fields[index].Precision = field.Precision

		fieldNameToValues[field.Field] = values
		colTypes = append(colTypes, field.IsNumb)
	}

	// 生成指定数量行的数据
	rows := make([][]string, 0)
	for i := 0; i < total; i++ {
		for _, field := range def.Fields {
			if !stringUtils.FindInArr(field.Field, fieldsToExport) {
				continue
			}

			values := fieldNameToValues[field.Field]
			fieldVal := values[i % len(values)]
			if len(rows) == i { rows = append(rows, make([]string, 0)) }
			rows[i] = append(rows[i], fieldVal)
		}
	}

	return rows, colTypes
}

func GenerateForField(field *model.DefField,  total int) []string {
	values := make([]string, 0)

	if len(field.Fields) > 0 { // sub fields
		arr := make([][]string, 0)
		for _, child := range field.Fields {
			childValues := GenerateForField(&child, total)
			arr = append(arr, childValues)
		}

		for i := 0; i < total; i++ {
			concat := ""
			for _, row := range arr {
				concat = concat + row[i]
			}

			concat = field.Prefix + concat + field.Postfix
			values = append(values, concat)
		}
		values = LoopSubFields(field, values, total)

	} else if field.From != "" && field.Select != "" { // refer to excel
		arr := strings.Split(field.From, ".")
		referField := constant.LoadedResValues[arr[0]]

		referField.From = field.From
		referField.Select = field.Select
		referField.Where = field.Where

		referField.Format = field.Format
		referField.Prefix = field.Prefix
		referField.Postfix = field.Postfix
		referField.Loop = field.Loop
		referField.Loopfix = field.Loopfix

		//values = GenerateFieldItemsFromDefinition(&referField, total)

	} else if field.From != "" && field.Range != "" { // refer to yaml file
		if field.Range != "" { // specific custom file
			//LoadDefinitionFromFile(constant.InputDir + field.Range)
		}

		//referField := constant.LoadedResValues[field.Field]
		//values = GenerateForField(&referField, total)
		// TODO: 需要处理range: small,large等逻辑

	} else {
		values = GenerateFieldItemsFromDefinition(field, total)
	}

	return values
}

func GenerateFieldItemsFromDefinition(field *model.DefField, total int) []string {
	if field.Loop == 0 {field.Loop = 1}

	values := make([]string, 0)

	// 整理出值的列表
	//datatype := strings.TrimSpace(field.Type)
	//if datatype == "" { datatype = "list" }

	fieldValue := GenerateList(field, total)

	index := 0
	count := 0
	for {
		// 处理格式、前后缀、loop等
		str := GenerateFieldValWithFix(*field, fieldValue, &index, true)
		values = append(values, str)

		count++
		if count >= total {
			break
		}
	}

	return values
}

func GenerateFieldVal(field model.DefField, fieldValue model.FieldValue, index *int) string {
	str := ""

	// 叶节点
	idx := *index % len(fieldValue.Values)
	val := fieldValue.Values[idx]
	str = GetFieldValStr(field, val)

	return str
}

func GetFieldValStr(field model.DefField, val interface{}) string {
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
		case string:
			str = val.(string)
			fmt.Sprintf(str)
		default:
	}

	return str
}

func LoopSubFields(field *model.DefField, oldValues []string, total int) []string {
	if field.Loop == 0 {field.Loop = 1}

	values := make([]string, 0)
	fieldValue := model.FieldValue{}

	for _, val := range oldValues {
		fieldValue.Values = append(fieldValue.Values, val)
	}

	index := 0
	count := 0
	for {
		// 处理格式、前后缀、loop等
		str := GenerateFieldValWithFix(*field, fieldValue, &index, false)
		values = append(values, str)

		count++
		if count >= total {
			break
		}
	}

	return values
}

func GenerateFieldValWithFix(field model.DefField, fieldValue model.FieldValue, indexOfRow *int, withLoop bool) string {
	prefix := field.Prefix
	postfix := field.Postfix

	loopStr := ""
	for j := 0; j < field.Loop; j++ {
		if loopStr != "" {
			loopStr = loopStr + field.Loopfix
		}

		str := GenerateFieldVal(field, fieldValue, indexOfRow)
		loopStr = loopStr + str

		*indexOfRow++
	}

	if withLoop {
		loopStr = prefix + loopStr + postfix
	}

	return loopStr
}